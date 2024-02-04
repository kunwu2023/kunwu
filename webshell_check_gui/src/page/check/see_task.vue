<template>
  <el-table :data="task_list.tableData" max-height="calc(100vh - 100px)" style="width: 100%">
    <el-table-column type="expand">
      <template #default="props">
          <taskResultList
              :task-id="props.row.id"
              :props-row="props.row"
              @open-new-window="openNewWindowCode"
          ></taskResultList>
      </template>
    </el-table-column>
    <el-table-column prop="dirPath" label="任务名" sortable/>
    <el-table-column prop="status" label="状态" width="100" sortable>
      <template #default="scope">
        <el-button v-if="scope.row.status === 1" size="small" type="info" :loading-icon="UploadFilled" loading plain>准备中</el-button>
        <el-button v-if="scope.row.status === 2" size="small" type="warning" :icon="Loading" loading plain>检测中</el-button>
        <el-button v-if="scope.row.status === 3" size="small" type="primary" :icon="CircleCheck" plain>完成</el-button>
        <el-button v-if="scope.row.status === 4" size="small" type="danger" :icon="CircleClose" plain>ssh异常</el-button>
      </template>
    </el-table-column>
    <el-table-column prop="createdAt" label="任务创建时间" width="170" sortable>
      <template #default="scope">
        <span v-if="scope.row.createdAt">{{ formatTimestamp(scope.row.createdAt) }}</span>
      </template>
    </el-table-column>
    <el-table-column>
      <template #default="scope">
<!--        <el-button size="small" @click="openNewWindow(scope)" text><el-icon><View /></el-icon></el-button>-->
        <el-popconfirm @confirm="delTask(scope)" title="你确定要删除这个任务吗？">
          <template #reference>
            <el-button size="small" text><el-icon><Delete /></el-icon></el-button>
          </template>
        </el-popconfirm>
      </template>
    </el-table-column>
  </el-table>
  <!--抽屉内容-->
  <el-drawer
      v-model="drawer.drawerOpen"
      size="100%"
      :title="drawer.title"
  >
    <el-table :data="drawer.tableData" max-height="calc(100vh - 150px)" style="width: 100%">
      <el-table-column prop="path" label="文件名" show-overflow-tooltip sortable/>
      <el-table-column prop="results" label="结果" width="100" sortable>
        <template #default="scope">
          <el-button v-if="scope.row.results === '待检测'" size="small" type="info" :loading-icon="UploadFilled" loading plain>待检测</el-button>
          <el-button v-if="scope.row.results === '正常'" size="small" type="success" :icon="CircleCheck" plain>正常</el-button>
          <el-button v-if="scope.row.results === '恶意'" size="small" type="warning" :icon="Warning" plain>恶意</el-button>
          <el-button v-if="scope.row.results === '不支持'" size="small" type="danger" :icon="CircleClose" plain>不支持</el-button>
        </template>
      </el-table-column>
      <el-table-column prop="size" label="大小" width="100" sortable/>
      <el-table-column prop="modificationTime" label="修改时间" sortable>
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
            <el-button size="small" @click="openNewWindowCode(scope)" text><el-icon><View /></el-icon></el-button>
          </el-tooltip>
          <el-popconfirm @confirm="delFile(scope)" title="你确定要处置(删除)这个文件吗？">
            <template #reference>
              <el-button size="small" text><el-icon><Delete /></el-icon></el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>
  </el-drawer>
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
import {onMounted, onUnmounted, reactive} from "vue";
import {CircleCheck, CircleClose, UploadFilled, View, Warning, Loading, Delete} from '@element-plus/icons-vue'
import {ipcRenderer} from "electron";
import {base64Decode, formatTimestamp, isNewData} from "@/utils/utils";
import {Codemirror} from "vue-codemirror";
import {javascript} from "@codemirror/lang-javascript";
import {oneDark} from "@codemirror/theme-one-dark";
import {ElMessage} from "element-plus";
import taskResultList from "@/page/check/inner/taskResultList.vue";

export default {
  name: "see_task",
  components: {
    Delete,
    Codemirror,
    View,
    taskResultList
  },
  setup(){
    const dialogData = reactive({
      flag:false,
      extensions: [javascript(), oneDark],
      webshellData:"",
      title:""
    })
    const task_list = reactive({  // 任务列表
      tableData: []
    })
    const drawer = reactive({ // 抽屉数据
      drawerOpen: false,
      tableData: [],
      title: "",
    })
    const openNewWindowCode = async (scope, model) => {  // 打开代码查看窗口
      if (scope.row.size <= 100000){
        if (model.model === 2){
          // TODO 远程查看文件
          // 将内容发送给go程序，等待返回
          const response = await ipcRenderer.invoke('go-request', {
            function: 'getSshReadFile',
            data: scope.row.path,
            taskID: model.id
          });
          dialogData.title = scope.row.path
          dialogData.webshellData = base64Decode(response['fileData'])
          dialogData.flag = true
        }else{
          // 将内容发送给go程序，等待返回
          const response = await ipcRenderer.invoke('go-request', {
            function: 'getReadFile',
            data: scope.row.path,
          });
          dialogData.title = scope.row.path
          dialogData.webshellData = response['fileData']
          dialogData.flag = true
        }
      }else{
        ElMessage({
          message: '文件太大，无法浏览',
          type: 'warning',
        })
      }
    }
    const openNewWindow = async (scope) => {
      const quickCheckList = await ipcRenderer.invoke('go-request', {
        function: 'getTaskList',
        data: scope.row.id,
      });
      drawer.tableData = quickCheckList["taskBaseList"]
      drawer.title = `"${scope.row.dirPath}"检出结果如下`
      drawer.drawerOpen = true
    };
    let stopLoop = false
    onMounted(async () => {
      await getQuickCheckList()
      // console.log(stopLoop)
      while (!stopLoop) {
        await getQuickCheckList()
        await new Promise(resolve => setTimeout(resolve, 2000)); // 延迟2秒
      }
    })
    const getQuickCheckList = async () => {
      const quickCheckList = await ipcRenderer.invoke('go-request', {
        function: 'getTask'
      });
      // 模拟数据加载，延迟 2 秒
      setTimeout(() => {
        // 如果新数据与现有数据不同，则将新数据推送到 tableData 中
        if (isNewData(quickCheckList.taskBaseList, task_list.tableData)) {
          task_list.tableData = quickCheckList.taskBaseList
        }
      }, 1);
    }
    const delTask = async (scope) => {
      await ipcRenderer.invoke('go-request', {
        function: 'delTask',
        data: scope.row.id
      });
      ElMessage('删除任务成功')
    }
    onUnmounted(() => {
      stopLoop = true;
    })
    return {
      CircleClose, CircleCheck, Warning, UploadFilled,Loading,
      task_list,drawer,dialogData,
      formatTimestamp,openNewWindow,openNewWindowCode,delTask
    }
  }
}
</script>

<style scoped>

</style>
