const { app, BrowserWindow, ipcMain, dialog } = require("electron");

let win;

function createWindow() {
  const win = new BrowserWindow({
    width: 1200,
    height: 600,
    webPreferences: {
      nodeIntegration: true,
      javascript: true,
      webSecurity: false,
    },
  });

  process.env.NODE_ENV === "development"
    ? win.loadURL("http://127.0.0.1:3000")
    : win.loadFile("index.html");

  win.webContents.openDevTools();
}

app.whenReady().then(createWindow);

app.on("window-all-closed", () => {
  if (process.platform !== "darwin") {
    app.quit();
  }
});

app.on("activate", () => {
  if (BrowserWindow.getAllWindows().length === 0) {
    createWindow();
  }
});

ipcMain.on("select.directory", (ev) => {
  dialog
    .showOpenDialog(win, {
      title: "选择需要监视的文件夹",
      properties: ["openDirectory"],
    })
    .then((result) => ev.returnValue = result.filePaths).catch((e) => {
        console.log(e);
    });
});
