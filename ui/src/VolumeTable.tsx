import { Button, message, Table } from "antd";
import { ColumnsType } from "antd/lib/table";
import axios from "axios";
import React, { FC } from "react";
import { formatPercentage, formatSpace } from "./utils";

export interface VolumeTableProps {
  data: any;
  loading: boolean;
  onSelect?: (data: any) => void;
  removable?: boolean;
}

const VolumeTable: FC<VolumeTableProps> = ({
  data,
  loading,
  onSelect,
  removable,
}) => {
  const cols: ColumnsType<any> = [
    {
      title: "卷名",
      dataIndex: "volume",
      key: "volume",
    },
    {
      title: "总空间",
      dataIndex: "totalSpace",
      key: "totalSpace",
      render: (_, r) =>
        formatSpace(
          data?.last?.volumesSpace?.find((s: any) => s?.volume === r?.volume)
            ?.totalSpace
        ) ?? "待获取",
    },
    {
      title: "剩余空间",
      dataIndex: "freeSpace",
      key: "freeSpace",
      render: (_, r) =>
        formatSpace(
          data?.last?.volumesSpace?.find((s: any) => s?.volume === r?.volume)
            ?.freeSpace
        ) ?? "待获取",
    },
    {
      title: "限额",
      dataIndex: "limit",
      key: "limit",
      render: (v, r) =>
        v !== 0
          ? r.limitPercentage
            ? formatPercentage(v)
            : formatSpace(v) ?? "未知"
          : "未限定",
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
              .delete("http://127.0.0.1:8005/volume", {
                params: {
                  volume: r.volume,
                },
              })
              .then(() => message.success("已移除" + r.volume))
              .catch((e) => `移除${r.volume}失败，${e}`);
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
      dataSource={data?.volumes}
      rowKey="volume"
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

export default VolumeTable;
