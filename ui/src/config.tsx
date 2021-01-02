import {
  BackTop,
  Button,
  InputNumber,
  message,
  Modal,
  Select,
  Space,
} from "antd";
import React, { FC } from "react";
import { SyncOutlined, UploadOutlined } from "@ant-design/icons";
import axios from "axios";
import { selectDirectory, useData } from "./utils";
import DirectoryTable from "./DirectoryTable";
import VolumeTable from "./VolumeTable";

export interface ConfigProps {}

const Config: FC<ConfigProps> = (props: ConfigProps) => {
  const { data, loading, run } = useData();

  return (
    <Space
      direction="vertical"
      style={{ width: "100%", background: "#fff", padding: "12px" }}
    >
      <BackTop visible>
        <Button type="primary" icon={<SyncOutlined />} onClick={run} />
      </BackTop>
      <Button
        icon={<UploadOutlined />}
        type="primary"
        onClick={() => {
          const volume = selectDirectory();
          if (!volume) {
            return;
          }
          let limit = 0;
          let unit = 1024 * 1024;
          let isPercentage = false
          Modal.info({
            title: "请输入限额",
            content: (
              <Space>
                <InputNumber onChange={(v) => v && (limit = v as number)} />
                <Select
                  defaultValue="mb"
                  style={{ width: 80 }}
                  onChange={(v: string) => {
                    switch (v) {
                      case "b":
                        unit = 1;
                        isPercentage = false;
                        break;
                      case "kb":
                        unit = 1024;
                        isPercentage = false;
                        break;
                      case "mb":
                        unit = 1024 * 1024;
                        isPercentage = false;
                        break;
                      case "gb":
                        unit = 1024 * 1024 * 1024;
                        isPercentage = false;
                        break;
                      case "tb":
                        unit = 1024 * 1024 * 1024 * 1024;
                        isPercentage = false;
                        break;
                      case "%":
                        unit = 1;
                        isPercentage = true;
                        break
                      default:
                        break;
                    }
                  }}
                >
                  <Select.Option value="%">%</Select.Option>
                  <Select.Option value="b">byte</Select.Option>
                  <Select.Option value="kb">KB</Select.Option>
                  <Select.Option value="mb">MB</Select.Option>
                  <Select.Option value="gb">GB</Select.Option>
                  <Select.Option value="tb">TB</Select.Option>
                </Select>
              </Space>
            ),
            async onOk() {
              try {
                await axios.put(
                  "http://127.0.0.1:8005/volume",
                  JSON.stringify({
                    volume,
                    limit: limit * unit,
                    limitPercentage: isPercentage,
                  }),
                  {
                    headers: {
                      "Content-Type": "application/json",
                    },
                  }
                );
                message.success(`已添加：${volume}`);
              } catch (e) {
                message.error(
                  `无法添加：${volume}，错误：${e?.response?.data}`
                );
              }
            },
          });
        }}
      >
        添加卷
      </Button>
      <VolumeTable data={data} loading={loading} removable />
      <Button
        icon={<UploadOutlined />}
        type="primary"
        onClick={() => {
          const directory = selectDirectory();
          if (!directory) {
            return;
          }
          let limit = 0;
          let unit = 1024 * 1024;
          Modal.info({
            title: "请输入限额",
            content: (
              <Space>
                <InputNumber onChange={(v) => v && (limit = v as number)} />
                <Select
                  defaultValue="mb"
                  style={{ width: 80 }}
                  onChange={(v: string) => {
                    switch (v) {
                      case "b":
                        unit = 1;
                        break;
                      case "kb":
                        unit = 1024;
                        break;
                      case "mb":
                        unit = 1024 * 1024;
                        break;
                      case "gb":
                        unit = 1024 * 1024 * 1024;
                        break;
                      case "tb":
                        unit = 1024 * 1024 * 1024 * 1024;
                        break;
                      default:
                        break;
                    }
                  }}
                >
                  <Select.Option value="b">byte</Select.Option>
                  <Select.Option value="kb">KB</Select.Option>
                  <Select.Option value="mb">MB</Select.Option>
                  <Select.Option value="gb">GB</Select.Option>
                  <Select.Option value="tb">TB</Select.Option>
                </Select>
              </Space>
            ),
            async onOk() {
              try {
                await axios.put(
                  "http://127.0.0.1:8005/directory",
                  JSON.stringify({
                    directory,
                    limit: limit * unit,
                  }),
                  {
                    headers: {
                      "Content-Type": "application/json",
                    },
                  }
                );
                message.success(`已添加：${directory}`);
              } catch (e) {
                message.error(
                  `无法添加：${directory}，错误：${e?.response?.data}`
                );
              }
            },
          });
        }}
      >
        添加文件夹
      </Button>
      <DirectoryTable data={data} loading={loading} removable />
    </Space>
  );
};

export default Config;
