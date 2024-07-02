import { Qa } from "../__generated__/graphql";
import {
  Button,
  Dialog,
  IconButton,
  Paper,
  Stack,
  SvgIcon,
  Table,
  TableBody,
  TableCell,
  TableRow,
  Typography,
} from "@mui/material";
import { truncateWords } from "../utils";
import { ReactComponent as TrashIcon } from "../assets/trash.svg";
import { ReactComponent as PencilIcon } from "../assets/pencil.svg";
import { useMutation } from "@apollo/client";
import { useState } from "react";
import { gql } from "../__generated__";

const SelectedQA = ({
  selectedQA,
  onEdit,
}: {
  selectedQA: Omit<Qa, "disease">;
  onEdit: () => void;
}) => {
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);

  const [deleteQA, { loading: isDeletingQA }] = useMutation(
    gql(`mutation deleteQA($id: ID!) {
    deleteQA(input:{id:$id}) {
      _stub
    }
  }`),
    {
      variables: { id: selectedQA.id },
      refetchQueries: ["qas"],
    }
  );
  return (
    <>
      <Dialog
        open={isDeleteDialogOpen}
        onClose={() => setIsDeleteDialogOpen(false)}
        disableEscapeKeyDown={isDeletingQA}
      >
        <Stack padding={2} gap={4}>
          <Typography padding={2}>
            Are you sure you want to delete this QA entry?
          </Typography>
          <Stack direction="row" justifyContent="flex-end" gap={1}>
            <Button
              color="inherit"
              onClick={() => setIsDeleteDialogOpen(false)}
              disabled={isDeletingQA}
            >
              Cancel
            </Button>
            <Button
              color="error"
              variant="contained"
              onClick={() =>
                deleteQA().finally(() => setIsDeleteDialogOpen(false))
              }
              disabled={isDeletingQA}
            >
              {isDeletingQA ? "Loading" : "Delete"}
            </Button>
          </Stack>
        </Stack>
      </Dialog>
      <Stack
        component={Paper}
        sx={{
          minHeight: "100%",
          overflowX: "hidden",
          borderEndEndRadius: 0,
          borderStartEndRadius: 0,
        }}
        borderRadius={2}
        elevation={0}
      >
        <Stack direction="row" alignItems="center" padding={1}>
          <Typography
            fontWeight="bold"
            flexGrow={1}
            color="secondary"
            paddingX={1}
          >
            Article {selectedQA.article.id}
          </Typography>
          <IconButton color="warning" onClick={onEdit}>
            <SvgIcon component={PencilIcon} fontSize="small" />
          </IconButton>
          <IconButton color="error" onClick={() => setIsDeleteDialogOpen(true)}>
            <SvgIcon component={TrashIcon} fontSize="small" />
          </IconButton>
        </Stack>
        <Table
          sx={{
            "& .MuiTableCell-root": {
              "&:first-of-type": { width: 120 },
            },
          }}
        >
          <TableBody>
            <TableRow>
              <TableCell>
                <Typography variant="body2" fontWeight="bold">
                  Article Text
                </Typography>
              </TableCell>
              <TableCell>
                {truncateWords(selectedQA.article.text, 50)}
              </TableCell>
            </TableRow>
            {selectedQA.questions.map((question, i) => (
              <TableRow key={question.id}>
                <TableCell>
                  <Typography variant="body2" fontWeight="bold">
                    Question {i + 1}
                  </Typography>
                </TableCell>
                <TableCell>
                  <Stack direction="row" gap={1}>
                    <Typography variant="body2" fontWeight="bold">
                      Text:
                    </Typography>
                    <Typography variant="body2">{question.text}</Typography>
                  </Stack>
                  {question.answers.map((answer, j) => (
                    <Stack key={j}>
                      <Stack direction="row" gap={1}>
                        <Typography variant="body2" fontWeight="bold">
                          Answer:
                        </Typography>
                        <Typography variant="body2">{answer.text}</Typography>
                      </Stack>
                      <Stack direction="row" gap={1}>
                        <Typography variant="body2" fontWeight="bold">
                          Answer Start:
                        </Typography>
                        <Typography variant="body2">
                          {answer.answer_start}
                        </Typography>
                      </Stack>
                    </Stack>
                  ))}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </Stack>
    </>
  );
};

export default SelectedQA;
