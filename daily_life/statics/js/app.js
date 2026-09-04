const scheduleForm = document.getElementById("scheduleForm");

const scheduleList = document.getElementById("scheduleList");


scheduleForm.addEventListener("submit", function(event) {

    event.preventDefault();


    const time = document.getElementById("time").value;

    const activity = document.getElementById("activity").value;


    const item = document.createElement("div");

    item.classList.add("task");


    item.innerHTML = `
        <div class="task-info">

            <span>
                ${time} - ${activity}
            </span>

        </div>

        <button class="delete-schedule">
            Delete
        </button>
    `;


    const emptyMessage = scheduleList.querySelector(".empty");

    if (emptyMessage) {
        emptyMessage.remove();
    }


    scheduleList.appendChild(item);


    document.getElementById("time").value = "";

    document.getElementById("activity").value = "";

});