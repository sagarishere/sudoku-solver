import React, { useState } from 'react';
import PropTypes from 'prop-types';
import boards from './boardsData';
import './App.css';

function App() {
  const [originalBoard, setOriginalBoard] = useState([]);
  const [solvedBoard, setSolvedBoard] = useState([]);

  const handleBoardChange = (index) => {
    setOriginalBoard(boards[index].slice()); // Make sure to slice to copy the board, not reference
    setSolvedBoard([]);
  };

  const solveBoard = async () => {
    try {
      const response = await fetch('http://localhost:8080/solve', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ board: originalBoard }),
      });

      if (!response.ok) {
        throw new Error('Network response was not ok ' + response.statusText);
      }

      const data = await response.json();
      setSolvedBoard(data.board);
    } catch (error) {
      console.error("Failed to fetch: ", error);
      alert("There was an error processing your request: " + error.message);
    }
  };


  return (
    <div className="App">
      <h1>Sudoku Solver</h1>
      <select onChange={(e) => handleBoardChange(e.target.value)} defaultValue="">
        <option value="" disabled>Select a board</option>
        {boards.map((_, index) => (
          <option key={`board-${index}`} value={index}>Board {index + 1}</option>
        ))}
      </select>
      <h2>Original Board</h2>
      <Board grid={originalBoard} solution={[]} />
      <button onClick={solveBoard}>Solve</button>
      <h2>Solved Board</h2>
      <Board grid={solvedBoard} solution={originalBoard} />
    </div>
  );
}

function Board({ grid, solution }) {
  return (
    <div className="centered-container">
      <table>
        <tbody>
          {grid.map((row, i) => (
            <tr key={`row-${i}`}>
              {row.map((cell, j) => (
                <td key={`${i}-${j}`} className={`cell${solution.length > 0 && solution[i][j] === 0 && cell !== 0 ? " solved" : ""}`}>
                  {cell === 0 ? '·' : cell}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

Board.propTypes = {
  grid: PropTypes.arrayOf(PropTypes.arrayOf(PropTypes.number.isRequired)).isRequired,
  solution: PropTypes.arrayOf(PropTypes.arrayOf(PropTypes.number.isRequired)), // Add propType for solution
};

export default App;
