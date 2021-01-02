import React, { useState } from 'react';
import { Layout, Menu } from "antd";
import Config from './config';
import Directory from './directory';
import Volume from './volume';

const { Header, Content, Footer } = Layout;

const App: React.FC = () => {
  const [page, setPage] = useState("directory");
  
  let content: React.ReactNode = "unknown page"

  switch (page) {
    case "config":
      content = <Config />;
      break
    case "directory":
      content = <Directory />;
      break
    case "volume":
      content = <Volume />;
      break
    default:
      break;
  }

  return (
    <Layout className="layout">
      <Header>
        <div className="logo" />
        <Menu
          theme="dark"
          mode="horizontal"
          selectedKeys={[page]}
          onSelect={(info) => setPage(info.key.toString())}
        >
          <Menu.Item key="config">配置</Menu.Item>
          <Menu.Item key="directory">文件夹</Menu.Item>
          <Menu.Item key="volume">卷</Menu.Item>
        </Menu>
      </Header>
      <Content
        style={{ padding: "20px 40px", minHeight: "calc(100vh - 136px)" }}
      >
        {content}
      </Content>
      <Footer style={{ textAlign: "center" }}>
        Space ©2020 Created by ronbb
      </Footer>
    </Layout>
  );
}

export default App;
