<template>
  <el-upload
      @dragover.prevent
      @drop="handleDropUpFile"
      @click="upFile"
      ref="uploadRef"
      class="upload-demo"
      drag
      :auto-upload="false"
      :disabled="true"
  >
    <el-icon class="el-icon--upload"><upload-filled /></el-icon>
    <div class="el-upload__text">
      {{ el_upload_text }} 或者<em>点击此处上传</em>
    </div>
    <template #tip>
    </template>
  </el-upload>
  <div class="el-upload__tip" v-if="check_result.tableData.length === 0">
    选择或拖拽入待检测的文件或文件夹即可开始检测，在这里的检测结果会在程序关闭后清空。
  </div>
  <div class="el-upload__tip" v-else>
    {{ check_result.result_txt }}
  </div>
  <el-table :data="check_result.tableData" max-height="calc(100vh - 350px)" style="width: 100%">
    <template #empty>
      <div class="el-table__empty-block">
        <span class="el-table__empty-text">{{ contentPrompt }}</span>
      </div>
    </template>
    <el-table-column prop="path" label="文件名" show-overflow-tooltip :sortable="sortableFlag"/>
    <el-table-column prop="results" label="结果" width="100" :sortable="sortableFlag">
      <template #default="scope">
        <el-button v-if="scope.row.results === '待检测'" size="small" type="info" :loading-icon="UploadFilled" loading plain>待检测</el-button>
        <el-button v-if="scope.row.results === '正常'" size="small" type="success" :icon="CircleCheck" plain>正常</el-button>
        <div v-if="scope.row.cloudResultsFlag !== 1">
          <el-button v-if="scope.row.results === '恶意'" size="small" type="warning" :icon="Cloudy" plain>恶意</el-button>
        </div>
        <el-button v-if="scope.row.results === '恶意' && scope.row.cloudResultsFlag === 1" size="small" type="warning" :icon="Warning" plain>恶意</el-button>
        <el-button v-if="scope.row.results === '不支持'" size="small" type="danger" :icon="CircleClose" plain>不支持</el-button>
      </template>
    </el-table-column>
    <el-table-column prop="size" label="大小" width="100" :sortable="sortableFlag">
      <template #default="scope">
        {{formatBytes(scope.row.size)}}
      </template>
    </el-table-column>
    <el-table-column prop="modificationTime" label="修改时间" :sortable="sortableFlag">
      <template #default="scope">
        {{ formatTimestamp(scope.row.modificationTime) }}
      </template>
    </el-table-column>
    <el-table-column label="操作">
      <template #default="scope">
        <el-tooltip
            class="box-item"
            effect="dark"
            content="查看内容"
            placement="top"
        >
        <el-button size="small" @click="openNewWindow(scope)" text><el-icon><View /></el-icon></el-button>
        </el-tooltip>
        <el-popconfirm @confirm="delFile(scope)" title="你确定要处置(删除)这个文件吗？">
          <template #reference>
            <el-button size="small" text><el-icon><Delete /></el-icon></el-button>
          </template>
        </el-popconfirm>
      </template>
    </el-table-column>
  </el-table>
  <!--代码查看窗口-->
  <el-dialog
      v-model="dialogData.flag"
      :title="dialogData.title"
      width="95%"
      draggable
      fullscreen
  >
    <codemirror
        v-model="dialogData.webshellData"
        placeholder="内容是空..."
        :autofocus="true"
        :indent-with-tab="true"
        :extensions="dialogData.extensions"
        :style="{ height: 'calc(100vh - 150px)', width:'100%' }"
    />
  </el-dialog>
</template>

<script>
import {reactive, ref, onMounted, onUnmounted} from "vue";
import {ipcRenderer} from "electron";
import {ElMessage, ElNotification} from 'element-plus'
import {CircleCheck, Warning, UploadFilled, CircleClose, View, Delete, Cloudy} from '@element-plus/icons-vue'
import fs from 'fs/promises';
import {formatBytes, formatTimestamp, isNewData} from "@/utils/utils";
import {javascript} from "@codemirror/lang-javascript";
import {oneDark} from "@codemirror/theme-one-dark";
import {Codemirror} from "vue-codemirror";
import {useStore} from "vuex";

export default {
  name: "quick_check",
  computed: {
    Cloudy() {
      return Cloudy
    }
  },
  components: {
    Codemirror,
    Delete,
    View,
    UploadFilled,
  },
  setup(){
    const contentPrompt = ref("未开始检测")
    const sortableFlag = ref(true)
    const store = useStore()
    const uploadRef = ref("")
    const el_upload_text = ref("将待检测的文件拖拽到此")
    const check_result = reactive({
      tableData: [],
      result_txt:"",
    })
    const data = reactive({
      webShellCode:"",
    })
    const dialogData = reactive({
      flag:false,
      extensions: [javascript(), oneDark],
      webshellData:"",
      title:""
    })
    let stopLoop = false;
    const upFile = async () => {
      const file = await ipcRenderer.invoke('dialog:openFile', "openFile")
      detectFiles([file])
    }
    async function handleDropUpFile(event) {
      // 拖拽上传
      event.preventDefault(); // 阻止默认的拖拽行为
      const files = event.dataTransfer.files; // 获取拖拽的文件列表

      // 处理文件列表，例如获取文件路径等操作
      const fileListPromises = Array.from(files).map(async (file) => {
        const filePath = file.path;
        const stats = await fs.stat(filePath);
        return {
          path: filePath,
          name: file.name,
          size: file.size,
          lastModified: file.lastModified,
          dirType: stats.isDirectory() ? 'directory' : 'file',
        };
      });
      const fileList = await Promise.all(fileListPromises);
      detectFiles(fileList);
    }
    const detectFiles = async (fileList) => {  // 检测恶意文件
      contentPrompt.value = "正在检测"
      ElNotification({
        title: "快速扫描",
        message: `开始扫描`,
        type: "success",
      })
      sortableFlag.value = false
      // 将内容发送给go程序，等待返回
      let cloudFlag
      if (store.state.scan_configuration.advanced.cloud_scan){
        cloudFlag = 1
      }else {
        cloudFlag = 2
      }
      await ipcRenderer.invoke('go-request', {
        function: 'quickCheckFile',
        data: fileList,
        cloudFlag: cloudFlag
      });
      console.log("we",fileList[0])
      if (fileList[0] !== null){
        ElNotification({
          title: "快速扫描",
          message: `检测完成结果如下`,
          type: "success",
        })
      }
      await getQuickCheckList()
      sortableFlag.value = true
      stopLoop = true;
      contentPrompt.value = "未检出到恶意样本"
    }
    function countFilesAndThreats(arr) {
      let totalFiles = arr.length;
      let suspiciousFiles = 0;

      for(let i = 0; i < totalFiles; i++) {
        if(arr[i].results === '恶意') {
          suspiciousFiles++;
        }
      }
      return `检测文件数：${totalFiles} 发现可疑文件：${suspiciousFiles}`;
    }

    const getQuickCheckList = async () => { // 获取快速检测列表
      try {
        const quickCheckList = await ipcRenderer.invoke('go-request', {
          function: 'quickCheckList'
        });
        // 如果新数据与现有数据不同，则将新数据推送到 tableData 中
        if (isNewData(quickCheckList.quickCheckList, check_result.tableData)) {
          check_result.tableData = quickCheckList.quickCheckList;
          check_result.result_txt = countFilesAndThreats(check_result.tableData)
        }
      } catch (error) {
        console.log("错误？",error)
      }
    }

    const openNewWindow = async (scope) => {  // 打开代码查看窗口
      if (scope.row.size <= 1000000){
        // 将内容发送给go程序，等待返回
        const response = await ipcRenderer.invoke('go-request', {
          function: 'getReadFile',
          data: scope.row.path
        });
        dialogData.title = scope.row.path
        dialogData.webshellData = response['fileData']
        dialogData.flag = true
      }else{
        ElMessage({
          message: '文件太大，无法浏览',
          type: 'warning',
        })
      }
    }
    const delFile = async (scope) => {
      // 将内容发送给go程序，等待返回
      const response = await ipcRenderer.invoke('go-request', {
        function: 'delFile',
        data: scope.row.path
      });
      ElNotification({
        title: "注意",
        message: `文件${response["msg"]}`,
        type: "warning",
      })
    }
    onMounted(async () => {
      await getQuickCheckList()
      // console.log(stopLoop)
      while (!stopLoop) {
        await new Promise(resolve => setTimeout(resolve, 2000)); // 延迟2秒
        await getQuickCheckList()
      }
    })
    onUnmounted(() => {
      stopLoop = true;
    })
    return{
      CircleClose, CircleCheck, Warning, UploadFilled,
      formatTimestamp,handleDropUpFile,upFile,openNewWindow,delFile,formatBytes,
      el_upload_text, uploadRef, check_result, data, dialogData, sortableFlag,contentPrompt
    }
}
}
</script>

<style scoped>

</style>