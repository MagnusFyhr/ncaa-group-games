import React from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { Box } from "@mui/material"; // Import Box from Material-UI
import "./App.css";
import Home from "./Home";
import Leaderboard from "./Leaderboard";
import PicksForm from "./PicksForm";
import Rules from "./Rules";
import NavigationBar from "./NavigationBar"; // Make sure to import NavigationBar

function App() {
  return (
    <Box sx={{ flexGrow: 1 }}>
      <BrowserRouter>
        <NavigationBar />
        <Routes>
          <Route path="/" exact element={<Home />} />
          <Route path="/leaderboard" element={<Leaderboard />} />
          <Route path="/picks" element={<PicksForm />} />
          <Route path="/rules" element={<Rules />} />
        </Routes>
      </BrowserRouter>
    </Box>
  );
}

export default App;
