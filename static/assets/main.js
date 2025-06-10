
//recordArrival js
    function recordArrival() {
        const userID = document.getElementById("userID").value;
        const eventID = document.getElementById("eventID").value;

        fetch("/arrive", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ userID: parseInt(userID), eventID: parseInt(eventID) })
        })
        .then(res => {
            if (!res.ok) throw new Error("Request failed");
            return res.json();
        })
        .then(data => {
            document.getElementById("arrivalResult").innerText =
            `Arrived at: ${data.arrivedAt}
            Distance: ${data.distanceKm.toFixed(3)} km
            Speed: ${data.speed.toFixed(3)} m/min`;
        })
        .catch(err => {
            document.getElementById("arrivalResult").innerText = "Error: " + err.message;
        });
        }
//create user js
    function createUser() {
    const username = document.getElementById("username").value;
    const email = document.getElementById("email").value;
    const latitudeDms = document.getElementById("latitudeDms").value;
    const longitudeDms = document.getElementById("longitudeDms").value;

    fetch("/users", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, email, latitudeDms, longitudeDms })
    })
    .then(res => {
      if (!res.ok) throw new Error("Request failed");
      return res.text(); // assuming plain response
    })
    .then(response => {
      document.getElementById("msg").innerText = "User created successfully!";
      document.getElementById("username").value = "";
      document.getElementById("email").value = "";
      document.getElementById("latitudeDms").value = "";
      document.getElementById("longitudeDms").value = "";
    })
    .catch(err => {
      document.getElementById("msg").innerText = "Error: " + err.message;
    });
  }  
// createEvent js
function releaseEvent() {
    const eventName = document.getElementById("eventName").value.trim();
    const releaseLatDMS = document.getElementById("releaseLatDMS").value.trim();
    const releaseLngDMS = document.getElementById("releaseLngDMS").value.trim();
    const responseEl = document.getElementById("releaseResponse");

    responseEl.style.display = "none";
    responseEl.textContent = "";

    fetch("/release", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ eventName, releaseLatDMS, releaseLngDMS })
    })
    .then(res => {
      if (!res.ok) return res.text().then(text => { throw new Error(text); });
      return res.json();
    })
    .then(data => {
      responseEl.style.display = "block";
      responseEl.style.color = "green";
      responseEl.textContent = `Success! Event ID: ${data.eventID}, Name: ${data.eventName}, Released at: ${new Date(data.releaseTime).toLocaleString()}`;
      document.getElementById("eventName").value = "";
      document.getElementById("releaseLatDMS").value = "";
      document.getElementById("releaseLngDMS").value = "";
    })
    .catch(err => {
      responseEl.style.display = "block";
      responseEl.style.color = "red";
      responseEl.textContent = "Error: " + err.message;
    });
  }