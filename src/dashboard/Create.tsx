import "../App.css";
import { useState } from "react";

export default function Create() {
  const [selectedSize, setSelectedSize] = useState<
    "banner" | "square" | "vertical"
  >("banner");
  const [imagePreview, setImagePreview] = useState<string | null>(null);
  const [imageName, setImageName] = useState<string>("");
  const [imageDimensions, setImageDimensions] = useState<{
    width: number;
    height: number;
  } | null>(null);
  const [levelId, setLevelId] = useState<string>("");
  const [levelValid, setLevelValid] = useState<boolean | null>(null);
  const [levelName, setLevelName] = useState<string>("");
  const [checkingLevel, setCheckingLevel] = useState<boolean>(false);

  const handleImageUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      setImageName(file.name);
      const reader = new FileReader();
      reader.onloadend = () => {
        setImagePreview(reader.result as string);
        const img = new window.Image();
        img.src = reader.result as string;
        img.onload = () => {
          setImageDimensions({ width: img.width, height: img.height });
        };
      };
      reader.readAsDataURL(file);
    }
  };

  const checkLevelValidity = async (id: string) => {
    if (!id || id.trim() === "") {
      setLevelValid(null);
      return;
    }

    setCheckingLevel(true);
    try {
      const formData = new URLSearchParams();
      formData.append("levelID", id);

      const response = await fetch("http://localhost:8081/api/proxy/level", {
        // probs change this to a relative path later
        method: "POST",
        headers: {
          "Content-Type": "application/x-www-form-urlencoded",
        },
        body: formData.toString(),
      });

      const data = await response.text();

      if (data === "-1" || data.trim() === "-1") {
        setLevelValid(false);
        setLevelName("");
      } else {
        setLevelValid(true);
        // (format: 1:128:2:Level Name:3:...)
        const parts = data.split(":");
        const levelNameIndex = parts.indexOf("2");
        if (levelNameIndex !== -1 && levelNameIndex + 1 < parts.length) {
          setLevelName(parts[levelNameIndex + 1]);
        } else {
          setLevelName("");
        }
      }
    } catch (error) {
      console.error("Error checking level validity:", error);
      setLevelValid(false);
    } finally {
      setCheckingLevel(false);
    }
  };

  const handleLevelIdChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const newLevelId = event.target.value;
    setLevelId(newLevelId);
    if (newLevelId.trim() === "") {
      setLevelValid(null);
      setLevelName("");
    }
  };

  function handleSubmit() {
    const img = document.createElement("img");
    img.src = imagePreview || "";
    img.onload = function () {
      let valid = false;
      if (selectedSize === "banner" && img.width === 1456 && img.height === 180)
        valid = true;
      if (
        selectedSize === "square" &&
        img.width === 1456 &&
        img.height === 1456
      )
        valid = true;
      if (
        selectedSize === "vertical" &&
        img.width === 180 &&
        img.height === 1456
      )
        valid = true;
      if (levelValid !== true) {
        alert("Please enter a valid Level ID.");
        return;
      }
      if (!valid) {
        alert(
          `Image dimensions do not match the selected size! Expected ${selectedSize} dimensions are: ${getExpectedDimensions(
            selectedSize
          )}`
        );
        return;
      }
      // submit loggin here, probs make it so it goes to the pending endpoint first then us to approve
      // the image should be processed the following {accountid}-{levelid}-{1-3}.png
      alert("Advertisement submitted!");
    };
    if (!imagePreview) {
      alert("Please select an image.");
    } else {
      img.dispatchEvent(new Event("load"));
    }
  }

  return (
    <>
      <h1 className="text-2xl font-bold mb-6">Create Advertisement</h1>
      <p className="text-lg mb-6">
        Select the size and upload an image for your advertisement.
      </p>
      <p className="text-lg mb-6">
        You can only upload one advertisement each size and it expires after 7
        days.
      </p>
      <div className="form-group mb-6">
        <label className="text-lg font-bold mb-2 block">
          Advertisement Size
        </label>
        <select
          className="custom-select"
          value={selectedSize}
          onChange={(e) =>
            setSelectedSize(e.target.value as "banner" | "square" | "vertical")
          }
        >
          <option value="banner">Banner (1456 x 180)</option>
          <option value="square">Square (1456 x 1456)</option>{" "}
          {/* ill figure out the ratio  for this later */}
          <option value="vertical">Vertical (180 x 1456)</option>
        </select>
      </div>
      <div className="form-group mb-6">
        <label className="text-lg font-bold mb-2 block">Upload Image</label>
        <div
          style={{
            display: "flex",
            flexDirection: "column",
            alignItems: "flex-start",
            gap: "16px",
          }}
        >
          <input
            type="file"
            id="image-upload"
            accept="image/*"
            onChange={handleImageUpload}
            style={{ display: "none" }}
          />
          <button
            className="nine-slice-button small"
            onClick={() => document.getElementById("image-upload")?.click()}
            style={{ alignSelf: "flex-start" }}
          >
            Import Image
          </button>
          <div
            style={{
              marginLeft: "8px",
              alignSelf: "flex-start",
              fontSize: "0.95em",
            }}
          >
            {imageName && (
              <span>
                <strong>{imageName}</strong>
                {imageDimensions && (
                  <span>
                    {" "}
                    &nbsp;({imageDimensions.width} x {imageDimensions.height})
                  </span>
                )}
              </span>
            )}
          </div>
          <div
            className="image-preview"
            style={{ display: "flex", justifyContent: "center", width: "100%" }}
          >
            <img
              src={imagePreview || "./assets/blacksquare.png"}
              alt="No Image Selected"
            />
          </div>
        </div>
      </div>
      <div className="form-group mb-6">
        <label className="text-lg font-bold mb-2 block">Level ID</label>
        <div style={{ display: "flex", alignItems: "center", gap: "1em" }}>
          <input
            className="custom-select"
            placeholder="Enter Level ID"
            value={levelId}
            onChange={handleLevelIdChange}
          />
          <button
            className="nine-slice-button small"
            type="button"
            style={{ padding: "0.5em 1em", fontSize: "1em" }}
            onClick={() => checkLevelValidity(levelId)}
            disabled={checkingLevel || !levelId.trim()}
          >
            Check
          </button>
        </div>
        <div style={{ marginTop: "8px", fontSize: "0.9em" }}>
          {checkingLevel && (
            <span style={{ color: "#888" }}>Checking level...</span>
          )}
          {!checkingLevel && levelValid === true && (
            <span style={{ color: "#4CAF50" }}>
              ✓ Valid level{levelName && `: ${levelName}`}
            </span>
          )}
          {!checkingLevel && levelValid === false && (
            <span style={{ color: "#f44336" }}>✗ Invalid level ID</span>
          )}
        </div>
      </div>
      <div
        className="form-group mb-6"
        style={{ display: "flex", justifyContent: "center" }}
      >
        <button className="nine-slice-button small" onClick={handleSubmit}>
          Submit
        </button>
      </div>
    </>
  );
}

function getExpectedDimensions(selectedSize: string) {
  switch (selectedSize) {
    case "banner":
      return "1456 x 180";
    case "square":
      return "1456 x 1456";
    case "vertical":
      return "180 x 1456";
    default:
      return "Unknown size";
  }
}
