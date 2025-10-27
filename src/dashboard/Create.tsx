import "../App.css";
import { useEffect, useState } from "react";

import WarningIcon from "@mui/icons-material/WarningOutlined";
import UploadFileIcon from "@mui/icons-material/UploadFileOutlined";
import SearchIcon from "@mui/icons-material/SearchOutlined";
import SyncIcon from "@mui/icons-material/SyncOutlined";
import CheckCircleIcon from "@mui/icons-material/CheckCircleOutlined";

export default function Create() {
  const [selectedSize, setSelectedSize] = useState<
    "banner" | "square" | "skyscraper"
  >("banner");
  const [imagePreview, setImagePreview] = useState<string | null>(null);
  const [imageName, setImageName] = useState<string>("");
  const [imageDimensions, setImageDimensions] = useState<{
    width: number;
    height: number;
  } | null>(null);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [userId, setUserId] = useState<string | null>(null);
  const [uploading, setUploading] = useState<boolean>(false);
  const [levelId, setLevelId] = useState<string>("");
  const [levelValid, setLevelValid] = useState<boolean | null>(null);
  const [levelName, setLevelName] = useState<string>("");
  const [checkingLevel, setCheckingLevel] = useState<boolean>(false);

  const handleImageUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      setSelectedFile(file);
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

  useEffect(() => {
    // fetch session to get logged-in user's id for filename
    fetch("/session", { credentials: "include" })
      .then((res) => (res.ok ? res.json() : null))
      .then((data) => {
        if (data?.id) setUserId(data.id);
      })
      .catch(() => setUserId(null));
  }, []);

  const checkLevelValidity = async (id: string) => {
    if (!id || id.trim() === "") {
      setLevelValid(null);
      return;
    }

    setCheckingLevel(true);
    try {
      const formData = new URLSearchParams();
      formData.append("level-id", id);

      const response = await fetch("/proxy/level", {
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

  async function handleSubmit() {
    if (!selectedFile || !imagePreview) {
      alert("Please select an image.");
      return;
    }

    if (levelValid !== true) {
      alert("Please enter a valid Level ID.");
      return;
    }

    if (!imageDimensions) {
      alert("Image dimensions unknown. Please re-import the image.");
      return;
    }

    // validate expected dimensions
    const dims = imageDimensions;
    let valid = false;
    if (selectedSize === "banner" && dims.width === 1456 && dims.height === 180)
      valid = true;
    if (
      selectedSize === "square" &&
      dims.width === 1456 &&
      dims.height === 1456
    )
      valid = true;
    if (
      selectedSize === "skyscraper" &&
      dims.width === 180 &&
      dims.height === 1456
    )
      valid = true;

    if (!valid) {
      alert(
        `Image dimensions do not match the selected size! Expected ${selectedSize} dimensions are: ${getExpectedDimensions(
          selectedSize
        )}`
      );
      return;
    }

    if (!userId) {
      alert("You must be logged in to submit an advertisement.");
      return;
    }

    setUploading(true);

    try {
      // convert to WebP using canvas
      const img = new Image();
      img.src = imagePreview;
      await new Promise<void>((resolve, reject) => {
        img.onload = () => resolve();
        img.onerror = () =>
          reject(new Error("Failed to load image for conversion"));
      });

      const canvas = document.createElement("canvas");
      canvas.width = img.width;
      canvas.height = img.height;
      const ctx = canvas.getContext("2d");
      if (!ctx) throw new Error("Canvas unsupported");
      ctx.drawImage(img, 0, 0);

      const webpBlob: Blob | null = await new Promise((resolve) =>
        canvas.toBlob((b) => resolve(b), "image/webp", 0.92)
      );

      let uploadBlob: Blob;
      let ext = "webp";
      if (webpBlob) {
        uploadBlob = webpBlob;
      } else {
        // fallback: use original file
        uploadBlob = selectedFile;
        ext = selectedFile.name.split(".").pop() || "png";
      }

      // const sizeNum = selectedSize === "banner" ? 1 : selectedSize === "square" ? 2 : 3;
      const filename = `${userId}.${ext}`;

      const formData = new FormData();
      formData.append("image-upload", uploadBlob, filename);
      formData.append("type", selectedSize);
      formData.append("level-id", levelId);

      const resp = await fetch("/ads/submit", {
        method: "POST",
        credentials: "include",
        body: formData,
      });

      if (resp.ok) {
        alert("Advertisement submitted!");
        // reset form
        setSelectedFile(null);
        setImagePreview(null);
        setImageName("");
        setImageDimensions(null);
        setLevelId("");
        setLevelValid(null);
        setLevelName("");

        const r = await resp.json();
        console.debug(`Ad of ID ${r["ad_id"]} stored at ${r["image_url"]}`);
      } else {
        const txt = await resp.text();
        console.error("Upload failed:", txt);
        alert("Failed to submit advertisement.");
      }
    } catch (err) {
      console.error(err);
      alert("An error occurred while processing the image.");
    } finally {
      setUploading(false);
    }
  }

  return (
    <>
      <h1 className="text-2xl font-bold mb-6">Create Advertisement</h1>
      <p className="text-lg">
        Select the size and upload an image for your advertisement.
      </p>
      <p className="text-sm mb-6 text-gray-500">
        You're only allowed to create 1 advertisement per type. Each
        advertisement expires after 7 days.{" "}
        <b>You may have a maximum of 3 active advertisements at a time.</b>{" "}
        Before it can be shown in game, your advertisement must first be
        approved by an admin.
      </p>
      {/* Rules */}
      <div className="text-sm text-orange-500 mb-6">
        <p>
          <WarningIcon /> Do not upload inappropriate or controversial
          advertisements.
        </p>
        <p>
          <WarningIcon /> Do not self-promote anything non-Geometry Dash related. Memes or well-known creators are allowed.
        </p>
        <p>
          <WarningIcon /> No profanity or excessive text in the advertisement.
        </p>
        <p>
          <WarningIcon /> Do not promote any harmful, illegal, or offensive
          material including both your Geometry Dash level and your
          advertisement!
        </p>
        <p>
          <WarningIcon /> AI Generated Advertisements are hard rejection. 
        </p>
        <p>
          <WarningIcon /> Violating this may result in a ban. <b>No appeals!</b>
        </p>
      </div>
      <div className="form-group mb-6">
        <label className="text-lg font-bold mb-2 block">
          Advertisement Size
        </label>
        <select
          className="custom-select"
          value={selectedSize}
          onChange={(e) =>
            setSelectedSize(
              e.target.value as "banner" | "square" | "skyscraper"
            )
          }
        >
          <option value="banner">Banner (1456 x 180)</option>
          <option value="square">Square (1456 x 1456)</option>{" "}
          {/* ill figure out the ratio  for this later */}
          <option value="skyscraper">Skyscraper (180 x 1456)</option>
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
            <UploadFileIcon /> Import Image
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
            <SearchIcon style={{ scale: 1.5 }} />
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
        <button
          className="nine-slice-button small"
          onClick={handleSubmit}
          disabled={uploading}
        >
          {uploading ? <SyncIcon /> : <CheckCircleIcon />}{" "}
          {uploading ? "Uploading..." : "Submit"}
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
    case "skyscraper":
      return "180 x 1456";
    default:
      return "Unknown size";
  }
}
