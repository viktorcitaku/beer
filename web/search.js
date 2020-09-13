const form = document.getElementById("search_form");

const beersTable = document.getElementById("beers");

const message = document.getElementById("message");

const searchBeer = () => {
  const formData = new FormData(form);

  let queryString = "";
  let count = 0;
  const criteriaCount = 3;
  for (let pair of formData.entries()) {
    ++count;

    if (pair[1] && pair[1] === "") {
      continue;
    }

    if (pair[0] === "pairing_food") {
      // special handling for pairing_food
      queryString += pair[0] + "=" + pair[1].replace(" ", "_");
    } else {
      queryString += pair[0] + "=" + pair[1];
    }

    if (count < criteriaCount) {
      queryString += "&";
    }
  }

  console.log(queryString);

  message.innerHTML = "Loading ...";

  // blocking call
  fetch(
      queryString && queryString.length > 0
          ? "/api/beers?" + queryString
          : "/api/beers",
      {
        method: "GET",
        mode: "cors", // no-cors, *cors, same-origin
        headers: {
          Accept: "application/json",
          // 'Content-Type': 'application/x-www-form-urlencoded',
        },
      }
  )
  .then((response) => {
    if (response.status !== 200) {
      throw new Error("Something went wrong!");
    } else {
      return response.json();
    }
  })
  .then((jsonData) => {
    populateTable(jsonData);

    message.innerHTML = "Done.";
  })
  .catch((err) => {
    message.innerHTML = "No data!";
    console.log(err);
  });

  return false;
};

const populateTable = (jsonData) => {
  let data = "";

  if (jsonData.length > 0) {
    for (let i = 0; i < jsonData.length; i++) {
      data += "<tr>";
      data += "<td>" + jsonData[i].id + "</td>";
      data += "<td>" + jsonData[i].name + "</td>";
      data +=
          '<td><button onclick="saveBeer(' + jsonData[i].id + ', \'' + jsonData[i].name +'\'' + ')">Save beer</button></td>';
      data += "</tr>";
    }
  }

  beersTable.innerHTML = data;
};

const saveBeer = (beerId, beerName) => {

  const payload = {
    id: beerId,
    name: beerName
  }

  fetch('/api/save-beer',
      {
        method: "POST",
        mode: "cors", // no-cors, *cors, same-origin
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload)
      }
  )
  .then((response) => {
    if (response.status !== 200) {
      throw new Error("Something went wrong!");
    } else {
      console.log(payload);
    }
  })
  .catch((err) => {
    message.innerHTML = "Unable to save beer with ID: " + beerId;
    console.log(err);
  });
};
