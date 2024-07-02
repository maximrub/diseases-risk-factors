/* eslint-disable */
import { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string;
  String: string;
  Boolean: boolean;
  Int: number;
  Float: number;
};

export type Answer = {
  __typename?: 'Answer';
  answer_start: Scalars['Int'];
  text: Scalars['String'];
};

export type AnswerInput = {
  answer_start: Scalars['Int'];
  text: Scalars['String'];
};

export type Article = {
  __typename?: 'Article';
  id: Scalars['ID'];
  text: Scalars['String'];
};

export type CreateQaInput = {
  articleId: Scalars['String'];
  diseaseId: Scalars['String'];
  questions?: InputMaybe<Array<QuestionInput>>;
};

export type DeleteQaInput = {
  id: Scalars['ID'];
};

export type DeleteQaPayload = {
  __typename?: 'DeleteQAPayload';
  _stub: Scalars['String'];
};

export type Disease = {
  __typename?: 'Disease';
  category: Scalars['String'];
  dbLinks: DiseaseDbLinks;
  description: Scalars['String'];
  id: Scalars['ID'];
  names: Array<Scalars['String']>;
};

export type DiseaseDbLinks = {
  __typename?: 'DiseaseDBLinks';
  icd10?: Maybe<Array<Scalars['String']>>;
  icd11?: Maybe<Array<Scalars['String']>>;
  mesh?: Maybe<Array<Scalars['String']>>;
};

export type Mutation = {
  __typename?: 'Mutation';
  createQA: QaPayload;
  deleteQA: DeleteQaPayload;
  fetchDiseases: Scalars['Boolean'];
  updateQA: UpdateQaPayload;
};


export type MutationCreateQaArgs = {
  input: CreateQaInput;
};


export type MutationDeleteQaArgs = {
  input: DeleteQaInput;
};


export type MutationUpdateQaArgs = {
  input: UpdateQaInput;
};

export type Qa = {
  __typename?: 'QA';
  article: Article;
  disease: Disease;
  id: Scalars['ID'];
  questions: Array<Question>;
};

export type QaPayload = {
  __typename?: 'QAPayload';
  qa?: Maybe<Qa>;
};

export type Query = {
  __typename?: 'Query';
  article: Article;
  diseases: Array<Disease>;
  qas?: Maybe<Array<Qa>>;
};


export type QueryArticleArgs = {
  id: Scalars['ID'];
};


export type QueryQasArgs = {
  diseaseId: Scalars['ID'];
};

export type Question = {
  __typename?: 'Question';
  answers: Array<Answer>;
  id: Scalars['ID'];
  text: Scalars['String'];
};

export type QuestionInput = {
  answers?: InputMaybe<Array<AnswerInput>>;
  text: Scalars['String'];
};

export type UpdateAnswerInput = {
  answer_start: Scalars['Int'];
  text: Scalars['String'];
};

export type UpdateQaInput = {
  patch: UpdateQaPatch;
  qaId: Scalars['ID'];
};

export type UpdateQaPatch = {
  questions?: InputMaybe<Array<UpdateQuestionInput>>;
};

export type UpdateQaPayload = {
  __typename?: 'UpdateQAPayload';
  qa?: Maybe<Qa>;
};

export type UpdateQuestionInput = {
  answers?: InputMaybe<Array<UpdateAnswerInput>>;
  text: Scalars['String'];
};

export type FetchDiseasesMutationVariables = Exact<{ [key: string]: never; }>;


export type FetchDiseasesMutation = { __typename?: 'Mutation', fetchDiseases: boolean };

export type DiseaseIdsQueryVariables = Exact<{ [key: string]: never; }>;


export type DiseaseIdsQuery = { __typename?: 'Query', diseases: Array<{ __typename?: 'Disease', id: string }> };

export type DiseasesQueryVariables = Exact<{ [key: string]: never; }>;


export type DiseasesQuery = { __typename?: 'Query', diseases: Array<{ __typename?: 'Disease', id: string, names: Array<string>, category: string, description: string, dbLinks: { __typename?: 'DiseaseDBLinks', icd10?: Array<string> | null, icd11?: Array<string> | null, mesh?: Array<string> | null } }> };

export type ArticleQueryVariables = Exact<{
  articleId: Scalars['ID'];
}>;


export type ArticleQuery = { __typename?: 'Query', article: { __typename?: 'Article', id: string, text: string } };

export type CreateQaMutationVariables = Exact<{
  diseaseId: Scalars['String'];
  articleId: Scalars['String'];
  questions: Array<QuestionInput> | QuestionInput;
}>;


export type CreateQaMutation = { __typename?: 'Mutation', createQA: { __typename?: 'QAPayload', qa?: { __typename?: 'QA', id: string } | null } };

export type UpdateQaMutationVariables = Exact<{
  id: Scalars['ID'];
  questions: Array<UpdateQuestionInput> | UpdateQuestionInput;
}>;


export type UpdateQaMutation = { __typename?: 'Mutation', updateQA: { __typename?: 'UpdateQAPayload', qa?: { __typename?: 'QA', id: string } | null } };

export type QasQueryVariables = Exact<{
  diseaseId: Scalars['ID'];
}>;


export type QasQuery = { __typename?: 'Query', qas?: Array<{ __typename?: 'QA', id: string, article: { __typename?: 'Article', id: string, text: string }, questions: Array<{ __typename?: 'Question', id: string, text: string, answers: Array<{ __typename?: 'Answer', answer_start: number, text: string }> }> }> | null };

export type DeleteQaMutationVariables = Exact<{
  id: Scalars['ID'];
}>;


export type DeleteQaMutation = { __typename?: 'Mutation', deleteQA: { __typename?: 'DeleteQAPayload', _stub: string } };


export const FetchDiseasesDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"fetchDiseases"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"fetchDiseases"}}]}}]} as unknown as DocumentNode<FetchDiseasesMutation, FetchDiseasesMutationVariables>;
export const DiseaseIdsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"diseaseIds"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"diseases"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}}]}}]} as unknown as DocumentNode<DiseaseIdsQuery, DiseaseIdsQueryVariables>;
export const DiseasesDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"diseases"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"diseases"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"names"}},{"kind":"Field","name":{"kind":"Name","value":"dbLinks"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"icd10"}},{"kind":"Field","name":{"kind":"Name","value":"icd11"}},{"kind":"Field","name":{"kind":"Name","value":"mesh"}}]}},{"kind":"Field","name":{"kind":"Name","value":"category"}},{"kind":"Field","name":{"kind":"Name","value":"description"}}]}}]}}]} as unknown as DocumentNode<DiseasesQuery, DiseasesQueryVariables>;
export const ArticleDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"article"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"articleId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"article"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"articleId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"text"}}]}}]}}]} as unknown as DocumentNode<ArticleQuery, ArticleQueryVariables>;
export const CreateQaDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"createQA"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"diseaseId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"articleId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"questions"}},"type":{"kind":"NonNullType","type":{"kind":"ListType","type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"QuestionInput"}}}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"createQA"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"ObjectValue","fields":[{"kind":"ObjectField","name":{"kind":"Name","value":"diseaseId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"diseaseId"}}},{"kind":"ObjectField","name":{"kind":"Name","value":"articleId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"articleId"}}},{"kind":"ObjectField","name":{"kind":"Name","value":"questions"},"value":{"kind":"Variable","name":{"kind":"Name","value":"questions"}}}]}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"qa"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}}]}}]}}]} as unknown as DocumentNode<CreateQaMutation, CreateQaMutationVariables>;
export const UpdateQaDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"updateQA"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"questions"}},"type":{"kind":"NonNullType","type":{"kind":"ListType","type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"UpdateQuestionInput"}}}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"updateQA"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"ObjectValue","fields":[{"kind":"ObjectField","name":{"kind":"Name","value":"qaId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}},{"kind":"ObjectField","name":{"kind":"Name","value":"patch"},"value":{"kind":"ObjectValue","fields":[{"kind":"ObjectField","name":{"kind":"Name","value":"questions"},"value":{"kind":"Variable","name":{"kind":"Name","value":"questions"}}}]}}]}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"qa"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}}]}}]}}]} as unknown as DocumentNode<UpdateQaMutation, UpdateQaMutationVariables>;
export const QasDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"qas"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"diseaseId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"qas"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"diseaseId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"diseaseId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"article"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"text"}}]}},{"kind":"Field","name":{"kind":"Name","value":"questions"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"text"}},{"kind":"Field","name":{"kind":"Name","value":"answers"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"answer_start"}},{"kind":"Field","name":{"kind":"Name","value":"text"}}]}}]}}]}}]}}]} as unknown as DocumentNode<QasQuery, QasQueryVariables>;
export const DeleteQaDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"deleteQA"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"deleteQA"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"ObjectValue","fields":[{"kind":"ObjectField","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}]}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"_stub"}}]}}]}}]} as unknown as DocumentNode<DeleteQaMutation, DeleteQaMutationVariables>;