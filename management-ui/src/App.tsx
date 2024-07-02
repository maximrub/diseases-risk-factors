import "./App.css";
import {
  Container,
  CssBaseline,
  LinearProgress,
  ThemeProvider,
  Toolbar,
} from "@mui/material";
import Diseases from "./components/Diseases";
import { ApolloClient, ApolloProvider, InMemoryCache } from "@apollo/client";
import { apiUrl, theme } from "./utils";
import { Box } from "@mui/system";
import Appbar from "./components/Appbar";
import { useAuth0 } from "@auth0/auth0-react";
import { useEffect, useState } from "react";

function App(): JSX.Element {
  const [jwt, setJwt] = useState<string>();
  const {
    loginWithRedirect,
    getAccessTokenSilently,
    isLoading,
    isAuthenticated,
  } = useAuth0();

  useEffect(() => {
    if (isLoading) {
      return;
    }
    if (!isAuthenticated) {
      loginWithRedirect();
      return;
    }
    getAccessTokenSilently().then(setJwt);
  }, [isLoading, isAuthenticated, loginWithRedirect, getAccessTokenSilently]);

  if (!jwt) {
    return (
      <ThemeProvider theme={theme}>
        <Container
          maxWidth="xs"
          sx={{
            height: "100%",
            justifyContent: "center",
            display: "flex",
            flexDirection: "column",
          }}
        >
          <LinearProgress
            variant="indeterminate"
            sx={{ borderRadius: 1, height: 8 }}
            color="success"
          />
        </Container>
      </ThemeProvider>
    );
  }
  const client = new ApolloClient({
    cache: new InMemoryCache(),
    uri: apiUrl,
    headers: {
      Authorization: `Bearer ${jwt}`,
    },
  });

  return (
    <ThemeProvider theme={theme}>
      <ApolloProvider client={client}>
        <CssBaseline />
        <Appbar />
        <Box height="calc(100% - 64px)">
          <Toolbar />
          <Diseases />
        </Box>
      </ApolloProvider>
    </ThemeProvider>
  );
}

export default App;
