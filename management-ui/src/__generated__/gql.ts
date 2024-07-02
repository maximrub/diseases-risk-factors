/* eslint-disable */
import * as types from './graphql';
import { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';

/**
 * Map of all GraphQL operations in the project.
 *
 * This map has several performance disadvantages:
 * 1. It is not tree-shakeable, so it will include all operations in the project.
 * 2. It is not minifiable, so the string of a GraphQL query will be multiple times inside the bundle.
 * 3. It does not support dead code elimination, so it will add unused operations.
 *
 * Therefore it is highly recommended to use the babel-plugin for production.
 */
const documents = {
    "mutation fetchDiseases {\n    fetchDiseases\n  }": types.FetchDiseasesDocument,
    "\n      query diseaseIds {\n        diseases {\n          id\n        }\n      }\n    ": types.DiseaseIdsDocument,
    "\n      query diseases {\n        diseases {\n          id\n          names\n          dbLinks {\n            icd10\n            icd11\n            mesh\n          }\n          category\n          description\n        }\n      }\n    ": types.DiseasesDocument,
    "query article($articleId:ID!) {\n    article(id:$articleId) {\n      id\n      text\n    }\n  }": types.ArticleDocument,
    "mutation createQA($diseaseId:String!, $articleId:String!, $questions:[QuestionInput!]!) {\n    createQA(input: {\n      diseaseId: $diseaseId\n      articleId: $articleId\n      questions: $questions\n    }) {\n      qa {\n        id\n      }\n    }\n  }": types.CreateQaDocument,
    "mutation updateQA($id:ID!, $questions:[UpdateQuestionInput!]!) {\n      updateQA(input: {\n        qaId: $id\n        patch: {\n          questions: $questions\n        }\n      }) {\n      qa {\n        id\n      }\n    }\n  }": types.UpdateQaDocument,
    "query qas($diseaseId:ID!) {\n      qas(diseaseId:$diseaseId) {\n        id\n        article {\n          id\n          text\n        }\n        questions {\n          id\n          text\n          answers {\n            answer_start\n            text\n          }\n        }\n      }\n    }": types.QasDocument,
    "mutation deleteQA($id: ID!) {\n    deleteQA(input:{id:$id}) {\n      _stub\n    }\n  }": types.DeleteQaDocument,
};

/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 *
 *
 * @example
 * ```ts
 * const query = gql(`query GetUser($id: ID!) { user(id: $id) { name } }`);
 * ```
 *
 * The query argument is unknown!
 * Please regenerate the types.
 */
export function gql(source: string): unknown;

/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "mutation fetchDiseases {\n    fetchDiseases\n  }"): (typeof documents)["mutation fetchDiseases {\n    fetchDiseases\n  }"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n      query diseaseIds {\n        diseases {\n          id\n        }\n      }\n    "): (typeof documents)["\n      query diseaseIds {\n        diseases {\n          id\n        }\n      }\n    "];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n      query diseases {\n        diseases {\n          id\n          names\n          dbLinks {\n            icd10\n            icd11\n            mesh\n          }\n          category\n          description\n        }\n      }\n    "): (typeof documents)["\n      query diseases {\n        diseases {\n          id\n          names\n          dbLinks {\n            icd10\n            icd11\n            mesh\n          }\n          category\n          description\n        }\n      }\n    "];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "query article($articleId:ID!) {\n    article(id:$articleId) {\n      id\n      text\n    }\n  }"): (typeof documents)["query article($articleId:ID!) {\n    article(id:$articleId) {\n      id\n      text\n    }\n  }"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "mutation createQA($diseaseId:String!, $articleId:String!, $questions:[QuestionInput!]!) {\n    createQA(input: {\n      diseaseId: $diseaseId\n      articleId: $articleId\n      questions: $questions\n    }) {\n      qa {\n        id\n      }\n    }\n  }"): (typeof documents)["mutation createQA($diseaseId:String!, $articleId:String!, $questions:[QuestionInput!]!) {\n    createQA(input: {\n      diseaseId: $diseaseId\n      articleId: $articleId\n      questions: $questions\n    }) {\n      qa {\n        id\n      }\n    }\n  }"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "mutation updateQA($id:ID!, $questions:[UpdateQuestionInput!]!) {\n      updateQA(input: {\n        qaId: $id\n        patch: {\n          questions: $questions\n        }\n      }) {\n      qa {\n        id\n      }\n    }\n  }"): (typeof documents)["mutation updateQA($id:ID!, $questions:[UpdateQuestionInput!]!) {\n      updateQA(input: {\n        qaId: $id\n        patch: {\n          questions: $questions\n        }\n      }) {\n      qa {\n        id\n      }\n    }\n  }"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "query qas($diseaseId:ID!) {\n      qas(diseaseId:$diseaseId) {\n        id\n        article {\n          id\n          text\n        }\n        questions {\n          id\n          text\n          answers {\n            answer_start\n            text\n          }\n        }\n      }\n    }"): (typeof documents)["query qas($diseaseId:ID!) {\n      qas(diseaseId:$diseaseId) {\n        id\n        article {\n          id\n          text\n        }\n        questions {\n          id\n          text\n          answers {\n            answer_start\n            text\n          }\n        }\n      }\n    }"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "mutation deleteQA($id: ID!) {\n    deleteQA(input:{id:$id}) {\n      _stub\n    }\n  }"): (typeof documents)["mutation deleteQA($id: ID!) {\n    deleteQA(input:{id:$id}) {\n      _stub\n    }\n  }"];

export function gql(source: string) {
  return (documents as any)[source] ?? {};
}

export type DocumentType<TDocumentNode extends DocumentNode<any, any>> = TDocumentNode extends DocumentNode<  infer TType,  any>  ? TType  : never;