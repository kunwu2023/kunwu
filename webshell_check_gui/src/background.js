'use strict'

import { app, protocol, BrowserWindow, dialog, ipcMain } from 'electron'
import { createProtocol } from 'vue-cli-plugin-electron-builder/lib'
import installExtension, { VUEJS3_DEVTOOLS } from 'electron-devtools-installer'
const isDevelopment = process.env.NODE_ENV !== 'production'
const fs = require('fs')


// Scheme must be registered before the app is ready
protocol.registerSchemesAsPrivileged([
  { scheme: 'app', privileges: { secure: true, standard: true } }
])

// TODO 创建一个go进程
const { spawn } = require('child_process');
const path = require('path');
const net = require('net');
const goPath = path.join(__dirname, '../bin', 'webshell_gui_go');
console.log("执行的二进制文件路径",goPath)
const goProcess = spawn(goPath);

async function getFileDetails(filePath) {
  const stats = await fs.promises.stat(filePath);
  const fileInfo = {
    lastModified: stats.mtimeMs,
    name: path.basename(filePath),
    path: filePath,
    size: stats.size,
    dirType: stats.isDirectory() ? 'directory' : 'file',
  };
  return fileInfo;
}

// 获取文件信息
async function handleFileOpen(flag) {
  const { canceled, filePaths } = await dialog.showOpenDialog({
    properties: [flag]
  });

  if (canceled) {
    console.log("File open dialog was cancelled");
    return null;
  } else {
    console.log("Selected file path: ", filePaths[0]);
    const fileInfo = await getFileDetails(filePaths[0]);
    console.log("File details: ", fileInfo);
    return fileInfo;
  }
}
async function createWindow() {
  // Create the browser window.
  const win = new BrowserWindow({
    width: 888,
    height: 666,
    title: '昆吾WebShell检测', // 添加该行
    webPreferences: {
      
      // Use pluginOptions.nodeIntegration, leave this alone
      // See nklayman.github.io/vue-cli-plugin-electron-builder/guide/security.html#node-integration for more info
      nodeIntegration: process.env.ELECTRON_NODE_INTEGRATION,
      contextIsolation: !process.env.ELECTRON_NODE_INTEGRATION
    },
  })

  if (process.env.WEBPACK_DEV_SERVER_URL) {
    // Load the url of the dev server if in development mode
    await win.loadURL(process.env.WEBPACK_DEV_SERVER_URL)
    if (!process.env.IS_TEST) win.webContents.openDevTools()
  } else {
    createProtocol('app')
    // Load the index.html when not in development
    win.loadURL('app://./index.html')
  }
}

// Quit when all windows are closed.
app.on('window-all-closed', () => {
  // On macOS it is common for applications and their menu bar
  // to stay active until the user quits explicitly with Cmd + Q
  if (process.platform !== 'darwin') {
    app.quit()
  }
})

// 添加这个事件监听器，用于在应用退出之前结束子进程
app.on('before-quit', () => {
  goProcess.kill();
});

app.on('activate', () => {
  // On macOS it's common to re-create a window in the app when the
  // dock icon is clicked and there are no other windows open.
  if (BrowserWindow.getAllWindows().length === 0) createWindow()
})

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.on('ready', async () => {
  if (isDevelopment && !process.env.IS_TEST) {
    // Install Vue Devtools
    try {
      await installExtension(VUEJS3_DEVTOOLS)
    } catch (e) {
      console.error('Vue Devtools failed to install:', e.toString())
    }
  }
  createWindow()
  // TODO 自定义的IPC通讯算法
  ipcMain.handle('dialog:openFile', async (event, flag) => {
    return await handleFileOpen(flag);
  })  // 用ipcMain调用主进程的函数，为了后续渲染进程可以调用它
  // 下面是IPC通讯的部分
  let goPort;

  goProcess.stdout.on('data', (data) => {
    const match = data.toString().match(/Listening on .*:(\d+)/);
    if (match) {
      goPort = parseInt(match[1]);
      console.log(`Go process is listening on port ${goPort}`);
    } else {
      console.log(`Go process stdout: ${data}`);
    }
  });

  ipcMain.handle('go-request', async (event, request) => {
    console.log(" 向 Go 进程发送请求:", JSON.stringify(request));

    return new Promise((resolve, reject) => {
      const client = net.createConnection({ port: goPort }, () => {
        client.write(JSON.stringify(request));
      });

      let data = '';
      client.on('data', (chunk) => {
        data += chunk;
        // 检查数据是否以换行符结尾，表示数据已接收完整
        if (data.endsWith('\n')) {
          data = data.slice(0, -1); // 去掉换行符
          try {
            const response = JSON.parse(data);
            console.log(`Go process response: ${data}`);
            resolve(response);
          } catch (error) {
            console.log(`Go out not json: ${data}`);
            console.log('捕获到错误：', error.message);
            reject(error);
          }
          client.end();
        }
      });

      client.on('error', (err) => {
        console.error(`Go process connection error: ${err}`);
        reject(err);
      });
    });
  });

  goProcess.stdout.on('data', (data) => {
    console.error(`Go 表示不服，输出了: ${data}`);
  });
  goProcess.stderr.on('data', (data) => {
    console.error(`Go 表示报错了: ${data}`);
  });
})
app.name = '昆吾WebShell检测'
// Exit cleanly on request from parent process in development mode.
if (isDevelopment) {
  if (process.platform === 'win32') {
    process.on('message', (data) => {
      if (data === 'graceful-exit') {
        app.quit()
      }
    })
  } else {
    process.on('SIGTERM', () => {
      app.quit()
    })
  }
}