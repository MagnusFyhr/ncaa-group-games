// NavigationBar.js
import React from 'react';
import { AppBar, Toolbar, Typography, Button } from '@mui/material';
import { Link } from 'react-router-dom';

const NavigationBar = () => {
  return (
    <AppBar position="static" className="navbar">
      <Toolbar>
        <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
          The March Madness 16
        </Typography>
        <Button color="inherit" className="nav-button" component={Link} to="/">
          Home
        </Button>
        <Button color="inherit" className="nav-button" component={Link} to="/leaderboard">
          Leaderboard
        </Button>
        <Button color="inherit" className="nav-button" component={Link} to="/picks">
          Make/Edit Picks
        </Button>
        <Button color="inherit" className="nav-button" component={Link} to="/rules">
          Rules
        </Button>
      </Toolbar>
    </AppBar>
  );
};

export default NavigationBar;
