import { Button, message, Table } from "antd";
import { ColumnsType } from "antd/lib/table";
import axios from "axios";
import React, { FC } from "react";
import { formatPercentage, formatSpace } from "./utils";

export interface DirectoryTableProps {
  data: any;
  loading: boolean;
  onSelect?: (data: any) => void;
  removable?: boolean;
}

const DirectoryTable: FC<DirectoryTableProps> = ({ data, loading, onSelect, removable }) => {
  const cols: ColumnsType<any> = [
    {
      title: "文件夹",
      dataIndex: "directory",
      key: "directory",
    },
    {
      title: "大小",
      dataIndex: "usedSpace",
      key: "usedSpace",
      render: (_, r) =>
        formatSpace(
          data?.last?.directoriesSpace?.find(
            (s: any) => s?.directory === r?.directory
          )?.usedSpace
        ) ?? "待获取",
    },
    {
      title: "占比",
      dataIndex: "percentage",
      key: "percentage",
      render: (_, r) =>
        formatPercentage(
          data?.last?.directoriesSpace?.find(
            (s: any) => s?.directory === r?.directory
          )?.percentage
        ) ?? "待获取",
    },
    {
      title: "限额",
      dataIndex: "limit",
      key: "limit",
      render: (v) =>
      v!==0?
        formatSpace(
          v
        ) ?? "未知" : "未限定",
    },
  ];
  if (removable) {
    cols.push({
      title: "操作",
      dataIndex: "remove",
      key: "remove",
      render: (_, r) => (
        <Button
          size="small"
          onClick={() => {
            axios
              .delete("http://127.0.0.1:8005/directory", {
                params: {
                  directory: r.directory,
                },
              })
              .then(() => message.success("已移除" + r.directory))
              .catch((e) => `移除${r.directory}失败，${e}`);
          }}
        >
          移除
        </Button>
      ),
    });
  }
  return (
    <Table
      loading={loading}
      dataSource={data?.directories}
      rowKey="directory"
      rowSelection={
        onSelect && {
          type: "radio",
          onSelect,
        }
      }
      columns={cols}
    />
  );
};

export default DirectoryTable;
