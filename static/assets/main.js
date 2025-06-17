// === UI STATE ===
function toggleSidebar(forceClose = false) {
  const sidebar = document.getElementById("sidebar");
  if (forceClose) {
    sidebar.classList.add("collapsed");
  } else {
    sidebar.classList.toggle("collapsed");
  }
}

function setActiveLink(clickedLink) {
  const links = document.querySelectorAll(".sidebar a");
  links.forEach((link) => link.classList.remove("active"));
  if (clickedLink) clickedLink.classList.add("active");
}

// === CONTENT LOADING ===
function loadContent(file, element) {
  const mainContent = document.getElementById("mainContent");
  mainContent.innerHTML = "<p>Loading...</p>";

  fetch(file)
    .then((response) => {
      if (!response.ok) throw new Error("Failed to load " + file);
      return response.text();
    })
    .then((html) => {
      mainContent.innerHTML = html;
      setActiveLink(element);
      if (window.innerWidth <= 768) toggleSidebar(true);

      // Delay to ensure DOM is updated before calling fetchProfile()
      if (file === "profile.html") {
        setTimeout(fetchProfile, 0);
      }
    })
    .catch((err) => {
      console.error("[ERROR] Loading content:", err);
      mainContent.innerHTML = `<p style="color:red;">Error loading content: ${err.message}</p>`;
    });
}

// === PROFILE DATA FETCHING ===
function fetchProfile() {
  console.log("[DEBUG] Fetching profile...");

  fetch("/me")
    .then((res) => {
      console.log("[DEBUG] Response received:", res);
      if (!res.ok) throw new Error("Failed to fetch user profile");
      return res.json();
    })
    .then((user) => {
      console.log("[DEBUG] User profile data:", user);

      const fullName = `${user.firstname || ""} ${user.middlename || ""} ${
        user.lastname || ""
      }`.trim();

      // ONLY update profile page fields
      $("#profileUserName").text(user.username);
      $("#profileUserEmail").text(user.email);
      $("#profileFullname").text(fullName || "N/A");
    })
    .catch((err) => {
      console.error("[ERROR] Failed to load profile:", err);
      $("#profileUserName").text("N/A");
      $("#profileUserEmail").text("N/A");
      $("#profileFullname").text("N/A");
    });
}

// === BASIC USER DATA FOR SIDEBAR ===
function fetchUserInfo() {
  $.ajax({
    url: "/me",
    method: "GET",
    dataType: "json",
    success: function (user) {
      console.log("[INFO] User data fetched:", user);
      $("#userName").text(user.username);
      $("#userEmail").text(user.email);

      const fullName = `${user.firstname || ""} ${user.middlename || ""} ${
        user.lastname || ""
      }`.trim();
      $("#fullname").text(fullName || "N/A");
    },
    error: function (xhr) {
      console.warn("[WARN] Failed to fetch user info", xhr.status);
      $("#userName").text("Guest");
      $("#userEmail").text("N/A");
      $("#fullname").text("N/A");
    },
  });
}

// === LOGOUT ===
function logout() {
  console.log("[INFO] Logging out...");

  fetch("/logout", { method: "POST" })
    .then(() => {
      location.replace("login.html");
    })
    .catch((err) => {
      console.error("[ERROR] Logout failed:", err);
    });
}

// === AUTO LOAD DEFAULT PAGE ===
document.addEventListener("DOMContentLoaded", () => {
  fetchUserInfo();

  const defaultLink = document.querySelector(".sidebar a.active");
  if (defaultLink) {
    const match = defaultLink.getAttribute("onclick")?.match(/'([^']+)'/);
    if (match && match[1]) {
      loadContent(match[1], defaultLink);
    }
  }
});

// //create user js
async function createUser() {
  const username = document.getElementById("username").value.trim();
  const email = document.getElementById("email").value.trim();
  const password = document.getElementById("password").value.trim();
  const firstname = document.getElementById("firstname").value.trim();
  const middlename = document.getElementById("middlename").value.trim();
  const lastname = document.getElementById("lastname").value.trim();
  const latitudeDms = document.getElementById("latitudeDms").value.trim();
  const longitudeDms = document.getElementById("longitudeDms").value.trim();
  const msg = document.getElementById("msg");

  msg.className = "";
  msg.innerText = "";

  if (
    !username ||
    !email ||
    !password ||
    !firstname ||
    !middlename ||
    !lastname ||
    !latitudeDms ||
    !longitudeDms
  ) {
    msg.innerText = "Please fill in all fields.";
    msg.className = "error";
    return;
  }

  try {
    const response = await fetch("/users", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        username,
        email,
        password,
        firstname,
        middlename,
        lastname,
        latitudeDms,
        longitudeDms,
      }),
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(errorData.message || "Failed to create user.");
    }

    msg.innerText = "User created successfully!";
    msg.className = "success";

    // Clear form
    document.getElementById("username").value = "";
    document.getElementById("email").value = "";
    document.getElementById("password").value = "";
    document.getElementById("firstname").value = "";
    document.getElementById("middlename").value = "";
    document.getElementById("lastname").value = "";
    document.getElementById("latitudeDms").value = "";
    document.getElementById("longitudeDms").value = "";
  } catch (err) {
    msg.innerText = "Error: " + err.message;
    msg.className = "error";
  }
}

//recordArrival js
function recordArrival() {
  const userID = document.getElementById("userID").value;
  const eventID = document.getElementById("eventID").value;

  fetch("/arrive", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      userID: parseInt(userID),
      eventID: parseInt(eventID),
    }),
  })
    .then((res) => {
      if (!res.ok) throw new Error("Request failed");
      return res.json();
    })
    .then((data) => {
      document.getElementById("arrivalResult").innerText = `Arrived at: ${
        data.arrivedAt
      }
            Distance: ${data.distanceKm.toFixed(3)} km
            Speed: ${data.speed.toFixed(3)} m/min`;
    })
    .catch((err) => {
      document.getElementById("arrivalResult").innerText =
        "Error: " + err.message;
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
    body: JSON.stringify({ eventName, releaseLatDMS, releaseLngDMS }),
  })
    .then((res) => {
      if (!res.ok)
        return res.text().then((text) => {
          throw new Error(text);
        });
      return res.json();
    })
    .then((data) => {
      responseEl.style.display = "block";
      responseEl.style.color = "green";
      responseEl.textContent = `Success! Event ID: ${data.eventID}, Name: ${
        data.eventName
      }, Released at: ${new Date(data.releaseTime).toLocaleString()}`;
      document.getElementById("eventName").value = "";
      document.getElementById("releaseLatDMS").value = "";
      document.getElementById("releaseLngDMS").value = "";
    })
    .catch((err) => {
      responseEl.style.display = "block";
      responseEl.style.color = "red";
      responseEl.textContent = "Error: " + err.message;
    });
}

// //logout function
// function logout() {
//   fetch("/logout", {
//     method: "POST", // Or GET if your handler accepts it that way
//     credentials: "include",
//   })
//     .then((res) => {
//       if (!res.ok) throw new Error("Logout failed");
//       return res.json();
//     })
//     .then((data) => {
//       alert(data.message || "Logged out");
//       window.location.href = "/login.html";
//     })
//     .catch((err) => {
//       alert("Error logging out: " + err.message);
//     });
// }
