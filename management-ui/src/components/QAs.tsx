import { useQuery } from "@apollo/client";
import { Button, List, ListItemButton, Stack, Typography } from "@mui/material";
import { useState } from "react";
import { clamp } from "../utils";
import { gql } from "../__generated__";
import { Disease } from "../__generated__/graphql";
import QAForm from "./QAForm";
import SelectedQA from "./SelectedQA";

const QAs = ({ disease }: { disease: Disease }) => {
  const [selectedQAIndex, setSelectedQAIndex] = useState(0);
  const [qaFormMode, setQaFormMode] = useState<"closed" | "add" | "edit">(
    "closed"
  );
  const { data: qas } = useQuery(
    gql(`query qas($diseaseId:ID!) {
      qas(diseaseId:$diseaseId) {
        id
        article {
          id
          text
        }
        questions {
          id
          text
          answers {
            answer_start
            text
          }
        }
      }
    }`),
    { variables: { diseaseId: disease.id } }
  );
  const safeQAIndex = clamp(selectedQAIndex, 0, (qas?.qas?.length || 1) - 1);
  const selectedQA = qas?.qas?.[safeQAIndex];

  if (qaFormMode !== "closed") {
    const existingQA = qaFormMode === "edit" ? selectedQA : undefined;
    return (
      <QAForm
        onClose={() => setQaFormMode("closed")}
        diseaseId={disease.id}
        existingQA={existingQA}
      />
    );
  }

  if (!qas?.qas?.length) {
    return (
      <Stack alignItems="center" justifyContent="center" height="100%">
        <Typography variant="body2" color="text.secondary" margin={1}>
          Nothing found
        </Typography>
        <Button
          variant="contained"
          size="large"
          color="secondary"
          onClick={() => setQaFormMode("add")}
        >
          Add first QA
        </Button>
      </Stack>
    );
  }

  return (
    <Stack direction="row" gap={1} marginY={1} height="100%">
      <List
        sx={{
          paddingY: 0,
          "& .MuiListItemButton-root": {
            borderStartEndRadius: 8,
            borderEndEndRadius: 8,
            marginY: 1,
          },
        }}
      >
        {qas?.qas?.map((qa, i) => (
          <ListItemButton
            key={qa.id}
            selected={i === safeQAIndex}
            onClick={() => setSelectedQAIndex(i)}
          >
            Article {qa.article.id}
          </ListItemButton>
        ))}
        <ListItemButton onClick={() => setQaFormMode("add")}>
          <Typography variant="body2" color="primary">
            + Add New
          </Typography>
        </ListItemButton>
      </List>
      {selectedQA && (
        <SelectedQA
          selectedQA={selectedQA}
          onEdit={() => setQaFormMode("edit")}
        />
      )}
    </Stack>
  );
};

export default QAs;
