import "../App.css";

export default function Account() {
  return (
    <>
      <h1 className="text-2xl font-bold mb-6">My Account</h1>
      <p className="text-lg mb-6">
        Manage your account and view ads pending approval.
      </p>
      <div
        style={{
          display: "flex",
          gap: "1em",
          justifyContent: "center",
          alignItems: "center",
        }}
      >
        <button
          className="nine-slice-button"
          style={{ color: "#fff" }}
          onClick={async () => {
            if (
              !confirm(
                "Are you sure you want to delete your account? This cannot be undone."
              )
            )
              return;
            try {
              const res = await fetch("/api/account/delete", {
                method: "POST",
                credentials: "include",
              });
              if (res.ok) {
                alert("Account deleted. You will be logged out.");
                window.location.href = "/";
              } else {
                const txt = await res.text();
                console.error("Delete failed:", txt);
                alert("Failed to delete account.");
              }
            } catch (err) {
              console.error(err);
              alert("Failed to delete account.");
            }
          }}
        >
          Delete Account
        </button>

        <button
          className="nine-slice-button"
          onClick={async () => {
            try {
              const res = await fetch("/api/ads/get?status=pending", {
                method: "GET",
                credentials: "include",
              });
              if (res.ok) {
                const data = await res.json();
                // For now just open a new window with the JSON results for inspection
                const w = window.open();
                if (w) {
                  w.document.body.innerText = JSON.stringify(data, null, 2);
                } else {
                  alert("Pending ads received; check console.");
                  console.log(data);
                }
              } else {
                const txt = await res.text();
                console.warn("Pending ads endpoint:", res.status, txt);
                alert(
                  "Pending ads endpoint not implemented or returned an error."
                );
              }
            } catch (err) {
              console.error(err);
              alert("Failed to fetch pending ads.");
            }
          }}
        >
          Pending Ads
        </button>
      </div>
    </>
  );
}
