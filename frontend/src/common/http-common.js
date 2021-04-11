import axios from "axios";

export default axios.create({
    baseUR: "http://localhost:8081",
    headers: {
        "Content-type": "application/json"
    }

});