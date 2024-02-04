<template>
  <div style="position: relative;">
  <el-table :data="table.tableData" border stripe
            :empty-text="'未发现可疑文件'"
            max-height="calc(100vh - 250px)"
            v-loading="loadingFlag"
            element-loading-text="加载中..."
            :element-loading-spinner="svg"
            element-loading-svg-view-box="-10, -10, 50, 50"
            style="width: 100%"
            v-el-table-infinite-scroll="loadMore"
  >
    <el-table-column prop="path" label="文件名" show-overflow-tooltip sortable/>
    <el-table-column prop="results" label="结果" width="100" sortable>
      <template #default="scope">
        <div v-if="scope.row.cloudResultsFlag === 4">
          <el-button type="warning" size="small" :icon="Close">云端超时</el-button>
        </div>
        <div v-else>
          <el-button v-if="scope.row.results === '待检测'" size="small" type="info" :loading-icon="UploadFilled" loading plain>待检测</el-button>
          <el-button v-if="scope.row.results === '正常'" size="small" type="success" :icon="CircleCheck" plain>正常</el-button>
          <div v-if="scope.row.cloudResultsFlag !== 1">
            <el-button v-if="scope.row.results === '恶意'" size="small" type="warning" :icon="Cloudy" plain>恶意</el-button>
          </div>
          <el-button v-if="scope.row.results === '恶意' && scope.row.cloudResultsFlag === 1" size="small" type="warning" :icon="Warning" plain>恶意</el-button>
          <el-button v-if="scope.row.results === '不支持'" size="small" type="danger" :icon="CircleClose" plain>不支持</el-button>
        </div>
      </template>
    </el-table-column>
    <el-table-column prop="size" label="大小" width="100" sortable>
      <template #default="scope">
        {{formatBytes(scope.row.size)}}
      </template>
    </el-table-column>
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
        <el-popconfirm v-if="table.cloudModel" @confirm="delFile(scope)" title="你确定要处置(删除)这个文件吗？">
          <template #reference>
            <el-button size="small" text><el-icon><Delete /></el-icon></el-button>
          </template>
        </el-popconfirm>
      </template>
    </el-table-column>
  </el-table>
    <div v-if="table.noMoreData" style="text-align: center; margin-top: 10px;">没有更多数据了</div>
  </div>
</template>

<script>
import {ref, onMounted, reactive} from 'vue';
import { ipcRenderer } from 'electron';
import {formatBytes, formatTimestamp} from "@/utils/utils";
import {CircleCheck, CircleClose, Close, Cloudy, Delete, UploadFilled, View, Warning} from "@element-plus/icons-vue";
import {javascript} from "@codemirror/lang-javascript";
import {oneDark} from "@codemirror/theme-one-dark";
import {ElNotification} from "element-plus";
import elTableInfiniteScroll from 'el-table-infinite-scroll';

export default {
  computed: {
    Close() {
      return Close
    },
    Cloudy() {
      return Cloudy
    },
    CircleClose() {
      return CircleClose
    },
    Warning() {
      return Warning
    },
    CircleCheck() {
      return CircleCheck
    },
    UploadFilled() {
      return UploadFilled
    }
  },
  components: {View, Delete},
  directives: {
    'el-table-infinite-scroll': elTableInfiniteScroll
  },
  props: {
    taskId: Number,
    propsRow: null
  },
  setup(props, { emit }) {
    const loadingFlag = ref(true)

    const table = reactive({
      tableData:[],
      data:[],
      stepping:10,
      noMoreData:false,
      cloudModel:false,
    })

    const dialogData = reactive({
      flag:false,
      extensions: [javascript(), oneDark],
      webshellData:"",
      title:""
    })

    const loadMore = () =>{
      if (transferItems(table.data, table.tableData, table.stepping)) {
        table.noMoreData = true
      }else{
        table.noMoreData = false
      }
    }

    function transferItems(sourceArray, destinationArray, stepping) {
      if (sourceArray.length === 0){
        return true
      }else{
        const itemsToTransfer = sourceArray.splice(0, stepping);
        destinationArray.push(...itemsToTransfer);
        return false
      }
    }


    onMounted(async () => {
      loadingFlag.value = true
      const quickCheckList = await ipcRenderer.invoke('go-request', {
        function: 'getTaskList',
        data: props.taskId,
      });
      if (props.propsRow.model === 1) {
        table.cloudModel = true
      }else {
        table.cloudModel = false
      }
      // 模拟数据加载，延迟 2 秒
      table.data = quickCheckList["taskBaseList"];
      transferItems(table.data, table.tableData, table.stepping)
      loadingFlag.value = false;
    });

    const svg = `
        <path class="path" d="
          M 30 15
          L 28 17
          M 25.61 25.61
          A 15 15, 0, 0, 1, 15 30
          A 15 15, 0, 1, 1, 27.99 7.5
          L 15 15
        " style="stroke-width: 4px; fill: rgba(0, 0, 0, 0)"/>
      `

    const openNewWindowCode = async (scope) => {  // 打开代码查看窗口
      emit('open-new-window', scope, props.propsRow);
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
    return {
      formatTimestamp,openNewWindowCode,delFile,formatBytes,loadMore,
      table,dialogData,loadingFlag,svg
    };
  },
};
</script>

<style>
.item {
  margin-top: 10px;
  margin-right: 40px;
}
</style>
