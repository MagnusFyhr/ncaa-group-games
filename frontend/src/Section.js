import Paper from '@mui/material/Paper';
import { styled } from '@mui/material/styles';

export const Section = styled(Paper)(({ theme, height }) => ({
    height: height,
    padding: theme.spacing(2),
    ...theme.typography.body2,
    textAlign: 'center',
    elevation: 3
  }));