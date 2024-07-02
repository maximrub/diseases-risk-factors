import {
  Table,
  TableBody,
  TableCell,
  TableRow,
  Typography,
} from "@mui/material";
import { Disease } from "../__generated__/graphql";

const DiseaseInfo = ({ disease }: { disease: Disease }) => {
  return (
    <Table>
      <TableBody>
        <TableRow>
          <TableCell width={120}>
            <Typography variant="body2" fontWeight="bold">
              ID
            </Typography>
          </TableCell>
          <TableCell>{disease.id}</TableCell>
        </TableRow>
        <TableRow>
          <TableCell width={120}>
            <Typography variant="body2" fontWeight="bold">
              Names
            </Typography>
          </TableCell>
          <TableCell>{disease.names.join(", ")}</TableCell>
        </TableRow>
        <TableRow>
          <TableCell width={120}>
            <Typography variant="body2" fontWeight="bold">
              ICD10
            </Typography>
          </TableCell>
          <TableCell>{disease.dbLinks.icd10?.join(", ")}</TableCell>
        </TableRow>
        <TableRow>
          <TableCell width={120}>
            <Typography variant="body2" fontWeight="bold">
              ICD11
            </Typography>
          </TableCell>
          <TableCell>{disease.dbLinks.icd11?.join(", ")}</TableCell>
        </TableRow>
        <TableRow>
          <TableCell width={120}>
            <Typography variant="body2" fontWeight="bold">
              MESH
            </Typography>
          </TableCell>
          <TableCell>{disease.dbLinks.mesh?.join(", ")}</TableCell>
        </TableRow>
        <TableRow>
          <TableCell width={120}>
            <Typography variant="body2" fontWeight="bold">
              Category
            </Typography>
          </TableCell>
          <TableCell>{disease.category}</TableCell>
        </TableRow>
        <TableRow sx={{ "& *": { borderBottom: "none" } }}>
          <TableCell width={120}>
            <Typography variant="body2" fontWeight="bold">
              Description
            </Typography>
          </TableCell>
          <TableCell>{disease.description}</TableCell>
        </TableRow>
      </TableBody>
    </Table>
  );
};

export default DiseaseInfo;
