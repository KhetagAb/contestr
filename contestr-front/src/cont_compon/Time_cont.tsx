import { useState, useEffect } from "react";
// import type { GetRegattaContestStandingsResponses } from "../client";
import { useCurRegattaData } from "../data";


function Clock() {
  const {data} = useCurRegattaData();
  const [now, setNow] = useState(Date.now());

  useEffect(() => {
    const interval = setInterval(() => {
      setNow(Date.now());
    }, 1000);
    return () => clearInterval(interval);
  }, []);

  if (!data) {
    return null
  }

  const startTime = data.current_tour_start_time * 1000;
  // const endTime = startTime + data.current_tour_duration * 60 * 1000;
  const elapsed = Math.max(0, now - startTime);
  // const remaining = Math.max(0, endTime - now);

  const formatTime = (ms: number) => {
    const totalSeconds = Math.floor(ms / 1000);
    const hours = Math.floor(totalSeconds / 3600);
    const minutes = Math.floor((totalSeconds % 3600) / 60);
    const seconds = totalSeconds % 60;
    return `${hours.toString().padStart(2, "0")}:${minutes
        .toString()
        .padStart(2, "0")}:${seconds.toString().padStart(2, "0")}`;
  };

  return (
      // <div>
      //   <table>
      //     <tbody>
      //     <tr>
      //       <td>
      //         Время от начала тура: {formatTime(elapsed)}
      //       </td>
      //
      //     </tr>
      //     <tr>
      //       <td>
      //         Время до конца тура: {formatTime(remaining)}
      //       </td>
      //     </tr>
      //     </tbody>
      //   </table>
      //   <p></p>
      //   <p></p>
      // </div>

      <div>
        {formatTime(elapsed)}
        {/*До конца тура: {formatTime(remaining)}*/}
      </div>
  );
}

export default Clock;
