import { Box } from "@mui/material";
import type { ReactNode } from "react";
import TopBar from "./TopBar";
import SideBar from "./SideBar";

interface Props {
  children: ReactNode;
}

export default function AppShell({ children }: Props) {
  return (
    <Box sx={{ display: "flex", height: "100vh" }}>
      <TopBar />
      <SideBar />
      <Box
        component="main"
        sx={{
          flexGrow: 1,
          p: 2,
          mt: "64px", // height of AppBar
          ml: "240px", // width of Sidebar
        }}
      >
        {children}
      </Box>
    </Box>
  );
}
