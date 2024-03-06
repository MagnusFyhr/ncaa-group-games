import React, { useState, useEffect } from "react";
import { getPicksDeadline, getPointsBySeed } from './backend'; // Correct import statement

const Rules = () => {

    const [picksDeadline, setPicksDeadline] = useState('<insert_date>')
    const [pointsBySeed, setPointsBySeed] = useState('<insert_data>')
    useEffect(() => {
        const fetchPicksDeadline = async () => {
          try {
            const resp = await getPicksDeadline();
            console.log(resp)
            if (resp.error) {
                console.error('Error fetching picks deadline:', resp.error);
              return
            }
            setPicksDeadline(resp.picksDeadline)
          } catch (error) {
            console.error('Error fetching picks deadline:', error);
          }
        };
        const fetchPointsBySeed = async () => {
            try {
              const resp = await getPointsBySeed();
              console.log(resp)
              if (resp.error) {
                console.error('Error fetching points by seed:', resp.error);
                return
              }
            setPointsBySeed(resp.pointsPerSeed);
            } catch (error) {
              console.error('Error fetching points by seed:', error);
            }
          };
    
        fetchPicksDeadline();
        fetchPointsBySeed();
      }, []); // Empty dependency array ensures the effect runs only once after the initial render
    
    return (
        <div>
            <h1>How to Play 'The March Madness 16'!</h1>
            <h2>Making Picks</h2>
            <h5 >
              1. Head to MAKE/EDIT PICKS
            </h5>
            <h5>
              2. If you have not made picks yet<br/>
                - Fill out the form<br/>
                - enter your venmo name exactly as it appears in venmo<br/>
                - enter a PIN; this can be a super simple password that allows you, and only you, to edit your picks<br/>
                - DO NOT FORGET YOUR PIN! you will not be able to edit your picks.
                Please write it down, take a picture, or text it to yourself <br/>
                - Enter your first and last name, then get creative with your team name<br/>
                - Enter your picks and click submit!
            </h5>
            <h5>
              3. If you have already made picks <br/>
                - Enter your venmo name and PIN <br/>
                - Click the Load Picks button <br/>
                - Your info and picks should populate the screen
            </h5>
            <h2>Rules</h2>
            <h5>1. When creating picks your venmo name must match the venmo I received from, NO EXCEPTIONS.</h5>
            <h5>2. Each player is allowed to make ONE submission, NO EXCEPTIONS.</h5>
            <h5>3. After {picksDeadline} all picks are final, NO EXCEPTIONS.</h5>
            <h5>4. If venmos are not received before {picksDeadline} you will be refunded and picks voided, NO EXCEPTIONS.</h5>
            <h5>5. You CANNOT hedge by selecting two teams from the same game. Doing so will void the higher seeded pick.</h5>
            <h5>Example: you CANNOT choose Purdue and whichever 16 seed that they play. Doing so voids your Purdue pick</h5>
            <h5>6. Participant with most points at end of round of 64 (Thursday & Friday) wins 100% of pot.</h5>
            <h2>Payment & Prizes</h2>
            <h5>Venmo @magnus-fyhr $20 to be make your entry eligible</h5>
            <h5>On Sunday, the team with the most points as listed by leaderboard will receive 100% of the prize pool</h5>
            <h2>Scoring</h2>
            <h5>Each participant selects ONE team from every seed,<br/>
                a total of 16 selections from the Round of 64. <br/>
                ONE 1 seed, ONE 2 seed, ONE 3 seed and so on. <br/>
            </h5>
            <h5>
                Each correct pick from seeds selected will be scored as shown below:<br/>
                {Object.keys(pointsBySeed).map((key) => (
                    <React.Fragment key={key}>
                        {key} Seed : {pointsBySeed[key]} points earned<br/>
                    </React.Fragment>
                ))}
            </h5>
            <h5>The player with the most points at the end of the round of 64 wins</h5>
            <h5>In the event of a tiebreaker, player with most correct picks wins, </h5>
            <h5>followed by most correct lower seeded wins, ie. picking the 1, 2, & 3 seeds</h5>
            <h5>correctly would beat someone selecting 1,2, & 4 seeds correctly.</h5>
            <h5>In the unlikely event that is still a tie, the players split the pool</h5>
        </div>
    );
};

export default Rules;
