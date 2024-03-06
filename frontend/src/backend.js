export const getTeams = async () => {
    try {
      const response = await fetch('http://localhost:8080/teams');
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
      const resp = await response.json();
      return resp;
    } catch (error) {
      console.error('Error fetching teams:', error);
      return { error: error.message };
    }
};

export const getPointsBySeed = async () => {
  try {
    const response = await fetch('http://localhost:8080/points-per-seed');
    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }
    const resp = await response.json();
    return resp;
  } catch (error) {
    console.error('Error fetching points per seed:', error);
    return { error: error.message };
  }
};

export const getPicksDeadline = async () => {
  try {
    let response = await fetch('http://localhost:8080/deadline');
    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }
    const resp = await response.json();

    const utcDate = new Date(resp.picksDeadline);
    // console.log("UTC Date:", utcDate.toISOString());

    const offsetMinutes = utcDate.getTimezoneOffset() + 60;
    // console.log("Time Zone Offset (minutes):", offsetMinutes);

    const localTime = new Date(utcDate.getTime() - offsetMinutes * 60 * 1000);
    // console.log("Local Time:", localTime.toISOString());

    resp.picksDeadline = localTime.toISOString().slice(0,-5) + " CST";

    return resp;
  } catch (error) {
    console.error('Error fetching deadline:', error);
    return { error: error.message };
  }
}

export const getNumberOfPlayers = async () => {
  try {
    const response = await fetch('http://localhost:8080/number-of-players');
    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }
    const resp = await response.json();
    return resp;
  } catch (error) {
    console.error('Error fetching number of players:', error);
    return { error: error.message };
  }
};

export const getLeaderboard = async () => {
  try {
    const response = await fetch('http://localhost:8080/leaderboard');
    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }
    const resp = await response.json();
    return resp;
  } catch (error) {
    console.error('Error fetching teams:', error);
    return { error: error.message };
  }
};

export const putPicks = async (venmoName, pin, firstName, lastName, teamName, picks) => {
  const url = `http://localhost:8080/picks/${venmoName}`;
  const data = {
    pin,
    firstName,
    lastName,
    teamName,
    picks,
  };

  const body = JSON.stringify(data)
  try {
    console.log("broken?")
    const response = await fetch(url, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: body,
      timeout: 3000, // Set a timeout of 5 seconds
    });
    console.log("no?")
    if (response.status === 400) {
      return {error: response.message}
    }
    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }

    const resp = await response.json();
    return resp;
  } catch (error) {
    console.error('Error putting picks:', error);
    return { error: error.message };
  }
};

export const getPicks = async (venmoName, pin) => {
  const url = `http://localhost:8080/picks/${venmoName}/${pin}`;

  console.log(url)

  try {
    const response = await fetch(url)
    if (response.status === 404) {
      return {error: `No picks were found for ${venmoName}:${pin}`}
    }
    if (!response.ok) {
      return {error: `HTTP error! Status: ${response.status}`};
    }

    const resp = await response.json();
    return resp;
  } catch (error) {
    console.error(`Error getting picks for ${venmoName}:`, error);
    return { error: error.message };
  }
}