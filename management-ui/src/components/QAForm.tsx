import { useMutation, useQuery } from "@apollo/client";
import {
  Box,
  Button,
  Card,
  CardContent,
  IconButton,
  Paper,
  Snackbar,
  SvgIcon,
  TextField,
  Typography,
} from "@mui/material";
import { ReactComponent as CloseIcon } from "../assets/close.svg";
import { ReactComponent as TrashIcon } from "../assets/trash.svg";
import { Stack } from "@mui/system";
import { useState } from "react";
import { gql } from "../__generated__";
import { Qa, QuestionInput } from "../__generated__/graphql";

const QAForm = ({
  onClose,
  diseaseId,
  existingQA,
}: {
  onClose: () => void;
  diseaseId: string;
  existingQA?: Omit<Qa, "disease">;
}) => {
  const [articleId, setArticleId] = useState(
    existingQA ? existingQA.article.id : ""
  );
  const [questions, setQuestions] = useState<QuestionInput[]>(
    existingQA ? existingQA.questions.map(({ id, ...rest }) => rest) : []
  );
  const [chosenQuestionIndex, setChosenQuestionIndex] = useState<number>();
  const [submitFailed, setSubmitFailed] = useState(false);

  const { data: article } = useQuery(
    gql(`query article($articleId:ID!) {
    article(id:$articleId) {
      id
      text
    }
  }`),
    { variables: { articleId: articleId } }
  );

  const [createQA, { loading: creatingQA }] = useMutation(
    gql(`mutation createQA($diseaseId:String!, $articleId:String!, $questions:[QuestionInput!]!) {
    createQA(input: {
      diseaseId: $diseaseId
      articleId: $articleId
      questions: $questions
    }) {
      qa {
        id
      }
    }
  }`),
    {
      variables: { diseaseId, articleId, questions },
      onCompleted: () => onClose(),
      refetchQueries: ["qas"],
    }
  );

  const [updateQA, { loading: updatingQA }] = useMutation(
    gql(`mutation updateQA($id:ID!, $questions:[UpdateQuestionInput!]!) {
      updateQA(input: {
        qaId: $id
        patch: {
          questions: $questions
        }
      }) {
      qa {
        id
      }
    }
  }`),
    {
      variables: { id: existingQA?.id || "", questions },
      onCompleted: () => onClose(),
      refetchQueries: ["qas"],
    }
  );

  const isFormValid =
    questions.length &&
    questions.every(
      (question) =>
        question.text &&
        question.answers?.length &&
        question.answers?.every((answer) => answer.text)
    ) &&
    articleId;

  function appendAnswers() {
    const selection = window.getSelection();
    if (!selection || chosenQuestionIndex === undefined) return;
    const question = questions[chosenQuestionIndex];
    const answer = {
      text: selection.toString(),
      answer_start: selection.anchorOffset,
    };
    setQuestions((questions) => {
      const newQuestions = [...questions];
      newQuestions[chosenQuestionIndex] = {
        ...question,
        answers: [...(question.answers ?? []), answer],
      };
      return newQuestions;
    });
    setChosenQuestionIndex(undefined);
  }

  function appendQuestions() {
    setQuestions((questions) => [...questions, { text: "", answers: [] }]);
  }

  function setQuestionText(questionIndex: number, text: string) {
    const question = questions[questionIndex];
    setQuestions((questions) => {
      const newQuestions = [...questions];
      newQuestions[questionIndex] = {
        ...question,
        text,
      };
      return newQuestions;
    });
  }

  function removeAnswer(questionIndex: number, answerIndex: number) {
    const question = questions[questionIndex];
    setQuestions((questions) => {
      const newQuestions = [...questions];
      newQuestions[questionIndex] = {
        ...question,
        answers: question.answers?.filter((_, i) => i !== answerIndex),
      };
      return newQuestions;
    });
  }

  function submitForm() {
    if (creatingQA || updatingQA) {
      return;
    }
    if (!isFormValid) {
      setSubmitFailed(true);
      return;
    }
    if (existingQA) {
      return updateQA();
    }
    createQA();
  }

  return (
    <>
      <Box position="absolute" top={72} right={12}>
        <IconButton onClick={onClose} color="inherit">
          <SvgIcon component={CloseIcon} />
        </IconButton>
      </Box>
      <Stack paddingTop={3} direction="row" flexGrow={1}>
        <Stack
          gap={2}
          paddingX={1}
          flex={1}
          position="sticky"
          alignSelf="flex-start"
          top={80}
        >
          <TextField label="Disease ID" value={diseaseId} disabled fullWidth />
          <TextField
            label="Article ID"
            fullWidth
            value={articleId}
            onChange={(event) => setArticleId(event.target.value)}
            disabled={!!existingQA}
          />
        </Stack>
        <Box
          component={Paper}
          elevation={0}
          sx={{ borderStartStartRadius: 8, borderStartEndRadius: 8 }}
          borderRadius={0}
          flex={3}
          padding={2}
        >
          {article?.article?.text ? (
            <Typography
              variant="body2"
              whiteSpace="pre-wrap"
              onMouseUp={appendAnswers}
            >
              {article?.article?.text}
            </Typography>
          ) : (
            <Typography
              variant="body2"
              color="text.secondary"
              textAlign="center"
            >
              Article text will appear here
            </Typography>
          )}
        </Box>
        <Stack
          flex={2}
          paddingX={2}
          gap={2}
          paddingBottom={2}
          justifyContent="space-between"
        >
          {article?.article?.text && (
            <>
              <Stack gap={2}>
                {questions.map((question, i) => (
                  <Stack key={i} gap={1}>
                    <Stack
                      direction="row"
                      justifyContent="space-between"
                      alignItems="center"
                    >
                      <Typography fontWeight="bold" variant="body2">
                        Question {i + 1}
                      </Typography>
                      <IconButton
                        onClick={() =>
                          setQuestions((questions) =>
                            questions.filter((_, qi) => qi !== i)
                          )
                        }
                        color="error"
                        size="small"
                      >
                        <SvgIcon component={TrashIcon} fontSize="small" />
                      </IconButton>
                    </Stack>
                    <TextField
                      fullWidth
                      multiline
                      value={question.text}
                      onChange={(event) =>
                        setQuestionText(i, event.target.value)
                      }
                    />
                    {question.answers?.map((answer, j) => (
                      <Card
                        variant="outlined"
                        key={j}
                        sx={{ position: "relative" }}
                      >
                        <CardContent
                          sx={{
                            paddingTop: 1,
                            display: "flex",
                            gap: 1,
                            flexDirection: "column",
                          }}
                        >
                          <Typography
                            paddingY={0.5}
                            fontWeight="bold"
                            variant="body2"
                            color="primary"
                          >
                            Answer {j + 1}
                          </Typography>
                          <IconButton
                            onClick={() => removeAnswer(i, j)}
                            color="error"
                            size="small"
                            sx={{ position: "absolute", top: 4, right: 4 }}
                          >
                            <SvgIcon component={TrashIcon} fontSize="small" />
                          </IconButton>
                          <Typography variant="body2" whiteSpace="pre-wrap">
                            <b>Text:</b> {answer.text}
                          </Typography>
                          <Typography variant="body2">
                            <b>Answer start:</b> {answer.answer_start}
                          </Typography>
                        </CardContent>
                      </Card>
                    ))}
                    <Button
                      variant="outlined"
                      size="small"
                      onClick={() => setChosenQuestionIndex(i)}
                      disabled={chosenQuestionIndex === i}
                    >
                      {chosenQuestionIndex === i
                        ? "Please select an answer"
                        : "Mark answer"}
                    </Button>
                  </Stack>
                ))}
                <Button onClick={appendQuestions} color="secondary">
                  + Add question
                </Button>
              </Stack>
              <Button
                sx={{
                  position: "sticky",
                  bottom: 12,
                }}
                onClick={submitForm}
                variant="contained"
                color="success"
              >
                Submit
              </Button>
            </>
          )}
        </Stack>
      </Stack>
      <Snackbar
        open={submitFailed}
        autoHideDuration={3000}
        onClose={() => setSubmitFailed(false)}
        message="Failure. Please validate form"
        anchorOrigin={{
          vertical: "bottom",
          horizontal: "center",
        }}
      />
    </>
  );
};

export default QAForm;
