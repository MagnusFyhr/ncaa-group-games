// Bookmakers.js
import { useState, useEffect } from 'react';

import Grid from '@mui/material/Unstable_Grid2';
import Box from '@mui/material/Box';
import Stack from '@mui/material/Stack';
import { Typography } from '@mui/material';

import { Section } from './Section';
import CustomDataGrid from './CustomDataGrid';
import { getLeaderboard } from './backend';


function Leaderboard() {
  const [tableData, setTableData] = useState([]);
  const [error, setError] = useState('');


  const columns = [
    {
      "field": "teamName",
      "headerName": "Team Name",
      "width": 150
    },
    {
      "field": "firstName",
      "headerName": "First Name",
      "width": 100
    },
    {
      "field": "lastName",
      "headerName": "Last Name",
      "width": 100
    },
    {
      "field": "points",
      "headerName": "Points",
      "width": 100
    },
    {
      "field": "maxPoints",
      "headerName": "Ceiling Points",
      "width": 100
    },
    {
      "field": "picks",
      "headerName": "Picks",
      "width": 1000
    }
  ];

  useEffect(() => {
    const fetchLeaderboard = async () => {
      try {
        const resp = await getLeaderboard();
        console.log(resp);
        if (resp.error) {
          setError(resp.error)
          return
        }
        setTableData(resp.leaderboard)
      } catch (error) {
        console.error('Error fetching leaderboard:', error);
      }
    };

    fetchLeaderboard();
  }, []); // Empty dependency array ensures the effect runs only once after the initial render


  // Define a function to generate a unique id for each row
  const getRowId = (row) => row.id;

  return (
    <Box height='100%' sx={{ flexGrow: 1 }}>
        <h1>Leaderboard</h1>
        <Grid container spacing={2}>
            <Grid xs={12}>
            <Stack spacing={2}>
                <Section>
                  {tableData.length > 0 ? (
                    <CustomDataGrid rows={tableData} width="100%" columns={columns} getRowId={getRowId} />
                  ) : (
                    <Typography variant="body1">No data to display</Typography>
                  )}         
                  {error && <p style={{ color: 'red' }}>{error}</p>}        
                </Section>
            </Stack>
            </Grid>
        </Grid>
    </Box>
  );
}

export default Leaderboard;