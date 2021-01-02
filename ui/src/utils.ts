import { useRequest } from "ahooks";
import { message } from "antd";
import axios from "axios";

const { ipcRenderer } = window.require("electron");

export const selectDirectory = () => {
  const paths: string[] = ipcRenderer.sendSync("select.directory");
  if (!paths.length) {
    return;
  }
  return paths[0];
};

export const formatSpace = (s: number | undefined): string | undefined => {
  if (s === undefined) {
    return
  }
  if (s < 1024) {
    return s + "B";
  }
  s /= 1024;
  if (s < 1024) {
    return s.toFixed(2) + "KB";
  }
  s /= 1024;
  if (s < 1024) {
    return s.toFixed(2) + "MB";
  }
  s /= 1024;
  if (s < 1024) {
    return s.toFixed(2) + "GB";
  }
  s /= 1024;
  return s.toFixed(2) + "TB";
};

export const formatPercentage = (n: number | undefined) => {
  if (n === undefined) {
    return;
  }

  return n.toFixed(4) + "%";
};

export const useData = () => {
   return useRequest(
     async () => {
       const data: any = {};
       await Promise.all([
         axios
           .get("http://127.0.0.1:8005/volumes")
           .then((res) => (data.volumes = res.data)),
         axios
           .get("http://127.0.0.1:8005/directories")
           .then((res) => (data.directories = res.data)),
         axios
           .get("http://127.0.0.1:8005/last_record")
           .then((res) => (data.last = res.data)),
       ]).catch((e) => message.error(e?.response?.data));
       return data;
     },
     { pollingInterval: 5000, loadingDelay: 500 }
   );
}
