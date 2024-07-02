import { createTheme } from "@mui/material";

export const apiUrl = process.env.REACT_APP_API_URL;
export const auth0ClientId = process.env.REACT_APP_AUTH0_CLIENT_ID || "";
export const auth0Domain = process.env.REACT_APP_AUTH0_DOMAIN || "";
export const appUrl = process.env.REACT_APP_APP_URL;


export const theme = createTheme({
  palette: {
    mode: window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light",
  },
  components: {
    MuiPaper: {
      styleOverrides: {
        root: ({ theme }) => ({
          backgroundColor: theme.palette.grey[theme.palette.mode === "dark" ? 900 : 100],
        })
      }
    },
  },
  typography: {
    fontFamily: "Space Mono, monospace",
  }
});