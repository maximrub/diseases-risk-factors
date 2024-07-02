import { useMutation } from "@apollo/client";
import { useAuth0 } from "@auth0/auth0-react";
import {
  AppBar,
  Button,
  IconButton,
  Snackbar,
  SvgIcon,
  Toolbar,
  Tooltip,
  Typography,
} from "@mui/material";
import { useState } from "react";
import { gql } from "../__generated__";
import { ReactComponent as LogoutIcon } from "../assets/logout.svg";

const Appbar = () => {
  const [didSuccessfullyFetchDiseases, setDidSuccessfullyFetchDiseases] =
    useState<boolean | null>(null);

  const { logout, isAuthenticated } = useAuth0();

  const [fetchDiseases, { loading }] = useMutation(
    gql(`mutation fetchDiseases {
    fetchDiseases
  }`),
    {
      onCompleted: () => {
        setDidSuccessfullyFetchDiseases(true);
      },
      onError: () => {
        setDidSuccessfullyFetchDiseases(false);
      },
    }
  );

  return (
    <>
      <AppBar
        sx={{
          zIndex: (theme) => theme.zIndex.drawer + 1,
        }}
      >
        <Toolbar sx={{ gap: 2 }}>
          <Typography variant="h6" flexGrow={1}>
            Diseases Risk Factors Labler
          </Typography>
          <Button
            variant="outlined"
            color="inherit"
            disabled={loading}
            onClick={() => fetchDiseases()}
          >
            {loading ? "Loading" : "Fetch Diseases"}
          </Button>
          {isAuthenticated && (
            <Tooltip title="Log Out">
              <IconButton
                onClick={() =>
                  logout({ logoutParams: { returnTo: window.location.href } })
                }
                color="inherit"
              >
                <SvgIcon component={LogoutIcon} />
              </IconButton>
            </Tooltip>
          )}
        </Toolbar>
      </AppBar>
      <Snackbar
        open={didSuccessfullyFetchDiseases !== null}
        message={didSuccessfullyFetchDiseases ? "Success" : "Failure"}
        onClose={() => setDidSuccessfullyFetchDiseases(null)}
        autoHideDuration={2000}
      />
    </>
  );
};

export default Appbar;
