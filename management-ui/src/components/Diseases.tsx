import { useQuery } from "@apollo/client";
import {
  Divider,
  Drawer,
  InputBase,
  LinearProgress,
  ListItemButton,
  Stack,
  Tab,
  Tabs,
  Toolbar,
} from "@mui/material";
import { FC, useState } from "react";
import { gql } from "../__generated__";
import { FixedSizeList } from "react-window";
import { Container } from "@mui/system";
import DiseaseInfo from "./DiseaseInfo";
import QAs from "./QAs";

const Diseases: FC = () => {
  const [selectedItemIndex, setSelectedItemIndex] = useState<number>();
  const [selectedInnerPage, setSelectedInnerPage] = useState<"info" | "qas">(
    "info"
  );
  const [searchTerm, setSearchTerm] = useState<string>("");
  const { data: diseaseIds, loading } = useQuery(
    gql(`
      query diseaseIds {
        diseases {
          id
        }
      }
    `)
  );

  const { data: diseasesData } = useQuery(
    gql(`
      query diseases {
        diseases {
          id
          names
          dbLinks {
            icd10
            icd11
            mesh
          }
          category
          description
        }
      }
    `)
  );

  const foundDiseases = diseaseIds?.diseases
    .map((d, i) => ({ ...d, originalIndex: i }))
    .filter((d) => d.id.startsWith(searchTerm));

  const selectedItem = diseasesData?.diseases?.[selectedItemIndex ?? 0];

  if (loading) {
    return (
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
    );
  }

  return (
    <Stack direction="row" height="100%" width="100%">
      <Drawer
        variant="permanent"
        sx={{ width: 160, display: "flex", flexDirection: "column" }}
      >
        <Toolbar />
        <InputBase
          placeholder="Search by ID"
          onChange={(event) => setSearchTerm(event.target.value)}
          sx={{ width: 160, padding: 2 }}
        />
        <Divider />
        <FixedSizeList
          width={160}
          height={732}
          itemSize={40}
          itemCount={foundDiseases?.length ?? 0}
          overscanCount={50}
        >
          {(props) => {
            const { index, style } = props;
            const disease = foundDiseases?.[index];
            return (
              <ListItemButton
                style={style}
                sx={{
                  backgroundColor: (theme) => theme.palette.background.default,
                }}
                selected={selectedItem?.id === disease?.id}
                onClick={() => setSelectedItemIndex(disease?.originalIndex)}
              >
                {disease?.id}
              </ListItemButton>
            );
          }}
        </FixedSizeList>
      </Drawer>
      {selectedItem && (
        <Stack height="100%" width="100%">
          <Tabs value={selectedInnerPage === "info" ? 0 : 1}>
            <Tab label="Info" onClick={() => setSelectedInnerPage("info")} />
            <Tab
              label="QA Entries"
              onClick={() => setSelectedInnerPage("qas")}
            />
          </Tabs>
          {selectedInnerPage === "info" ? (
            <DiseaseInfo disease={selectedItem} />
          ) : (
            <QAs disease={selectedItem} />
          )}
        </Stack>
      )}
    </Stack>
  );
};

export default Diseases;
