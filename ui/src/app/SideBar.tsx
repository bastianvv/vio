import {
  Drawer,
  List,
  ListItemButton,
  ListItemText,
  Toolbar,
} from "@mui/material";
import { useLocation, useNavigate } from "react-router-dom";

const drawerWidth = 240;

export default function SideBar() {
  const navigate = useNavigate();
  const location = useLocation();

  return (
    <Drawer
      variant="permanent"
      sx={{
        width: drawerWidth,
        [`& .MuiDrawer-paper`]: {
          width: drawerWidth,
          boxSizing: "border-box",
        },
      }}
    >
      <Toolbar />
      <List>
        <ListItemButton
          selected={location.pathname.startsWith("/movies")}
          onClick={() => navigate("/movies")}
        >
          <ListItemText primary="Movies" />
        </ListItemButton>

        <ListItemButton
          selected={location.pathname.startsWith("/series")}
          onClick={() => navigate("/series")}
        >
          <ListItemText primary="Series" />
        </ListItemButton>
      </List>
    </Drawer>
  );
}
