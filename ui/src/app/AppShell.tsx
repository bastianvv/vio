import { Box } from "@mui/material";
import type { ReactNode } from "react";
import TopBar from "./TopBar";
import SideBar from "./SideBar";

interface Props {
  children: ReactNode;
}

export default function AppShell({ children }: Props) {
  return (
    <Box sx={{ display: "flex", minHeight: "100vh" }}>
      <TopBar />
      <SideBar />

      {/* DO NOT offset this manually */}
      <Box
        component="main"
        sx={{
          flexGrow: 1,
        }}
      >
        {children}
      </Box>
    </Box>
  );
}
