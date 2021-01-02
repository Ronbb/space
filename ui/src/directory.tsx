import { Line } from "@ant-design/charts";
import { SyncOutlined } from "@ant-design/icons";
import { BackTop, Button, DatePicker, message, Space } from "antd";
import axios from "axios";
import moment from "moment";
import React, { FC, useEffect, useState } from "react";
import DirectoryTable from "./DirectoryTable";
import { formatSpace, useData } from "./utils";

const { RangePicker } = DatePicker;

export interface DirectoryProps {}

const Directory: FC<DirectoryProps> = (props: DirectoryProps) => {
  const { data, run, loading } = useData();
  const [selectedDirectory, setSelectedDirectory] = useState<string>();
  const [range, setRange] = useState<moment.Moment[]>();
  const [space, setSpace] = useState<any>();

  useEffect(() => {
    if (!selectedDirectory || !range || range.length !== 2) {
      return;
    }
    axios
      .get("http://127.0.0.1:8005/space/directory", {
        params: {
          start: range[0].unix(),
          end: range[1].unix(),
          path: selectedDirectory,
        },
      })
      .then((res) => setSpace(res.data))
      .catch((e) =>
        message.error(
          `获取${selectedDirectory}的空间信息失败，${e?.response?.data}`
        )
      );
  }, [selectedDirectory, range]);

  return (
    <Space
      direction="vertical"
      style={{ width: "100%", background: "#fff", padding: "12px" }}
    >
      <BackTop visible>
        <Button type="primary" icon={<SyncOutlined />} onClick={run} />
      </BackTop>
      <RangePicker showTime onChange={(values) => setRange(values as any)} />
      <DirectoryTable
        data={data}
        loading={loading}
        onSelect={(e) => setSelectedDirectory(e?.directory)}
      />
      {space && (
        <Line
          data={space}
          xField="time"
          yField="usedSpace"
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
              text: "占用空间",
            },
            label: {
              formatter: (v) =>
                formatSpace(parseInt(v)),
            },
          }}
        />
      )}
    </Space>
  );
};

export default Directory;
