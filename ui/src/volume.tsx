import { Line } from "@ant-design/charts";
import { SyncOutlined } from "@ant-design/icons";
import { BackTop, Button, DatePicker, message, Space } from "antd";
import axios from "axios";
import moment from "moment";
import React, { FC, useEffect, useState } from "react";
import { useData, formatSpace, formatPercentage } from "./utils";
import VolumeTable from "./VolumeTable";

const { RangePicker } = DatePicker;

export interface VolumeProps {}

const Volume: FC<VolumeProps> = (props: VolumeProps) => {
  const { data, run, loading } = useData();
  const [selectedVolume, setSelectedVolume] = useState<string>();
  const [range, setRange] = useState<moment.Moment[]>();
  const [space, setSpace] = useState<any>();

  useEffect(() => {
    if (!selectedVolume || !range || range.length !== 2) {
      return;
    }
    axios
      .get("http://127.0.0.1:8005/space/volume", {
        params: {
          start: range[0].unix(),
          end: range[1].unix(),
          path: selectedVolume,
        },
      })
      .then((res) => setSpace(res.data))
      .catch((e) =>
        message.error(
          `获取${selectedVolume}的空间信息失败，${e?.response?.data}`
        )
      );
  }, [selectedVolume, range]);

  return (
    <Space
      direction="vertical"
      style={{ width: "100%", background: "#fff", padding: "12px" }}
    >
      <BackTop visible>
        <Button type="primary" icon={<SyncOutlined />} onClick={run} />
      </BackTop>
      <RangePicker showTime onChange={(values) => setRange(values as any)} />
      <VolumeTable
        data={data}
        loading={loading}
        onSelect={(e) => setSelectedVolume(e?.volume)}
      />
      {space && (
        <Line
          data={space}
          xField="time"
          yField="freeSpace"
          point={{
            size: 5,
            shape: "diamond",
            style: {
              fill: "white",
              stroke: "#2593fc",
              lineWidth: 2,
            },
          }}
          xAxis={{
            title: {
              text: "时间",
            },
            label: {
              formatter: (v) =>
                moment(parseInt(v) * 1000).format("YYYY-MM-DD HH:mm:ss"),
            },
          }}
          yAxis={{
            title: {
              text: "已用空间",
            },
            label: {
              formatter: (v, r) =>
                formatSpace(parseInt(v)) ??
                "" + formatPercentage(parseInt(v) / r.value.totalSpace) ??
                "",
            },
          }}
        />
      )}
    </Space>
  );
};

export default Volume;
