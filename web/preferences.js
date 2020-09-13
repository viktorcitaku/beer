let dateEditor = function (cell, onRendered, success, cancel, editorParams) {
  //cell - the cell component for the editable cell
  //onRendered - function to call when the editor has been rendered
  //success - function to call to pass the successfully updated value to Tabulator
  //cancel - function to call to abort the edit and return to a normal cell
  //editorParams - params object passed into the editorParams column definition property

  //create and style editor
  const editor = document.createElement("input");

  editor.setAttribute("type", "date");

  //create and style input
  editor.style.padding = "3px";
  editor.style.width = "100%";
  editor.style.boxSizing = "border-box";

  //Set value of editor to the current value of the cell
  editor.value = moment(cell.getValue(), "YYYY-MM-DD").format("YYYY-MM-DD")

  //set focus on the select box when the editor is selected (timeout allows for editor to be added to DOM)
  onRendered(function () {
    editor.focus();
    editor.style.css = "100%";
  });

  //when the value has been set, trigger the cell to update
  function successFunc() {
    success(moment(editor.value, "YYYY-MM-DD").format("YYYY-MM-DD"));
  }

  editor.addEventListener("change", successFunc);
  editor.addEventListener("blur", successFunc);

  //return the editor element
  return editor;
};

let uploadData = (data) => {
  console.log(data);

  fetch('/api/update-user-beer-preferences', {
        method: "POST",
        mode: "cors", // no-cors, *cors, same-origin
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data)
      }
  )
  .then((response) => {
    if (response.status !== 200) {
      throw new Error("Something went wrong!");
    } else {
      console.log("Data updated!")
    }
  })
  .catch((err) => {
    console.log(err);
  });
}

//create Tabulator on DOM element with id "example-table"
let table = new Tabulator("#beers", {
  height: "100%", // set height of table (in CSS or here), this enables the Virtual DOM and improves render speed dramatically (can be any valid css height value)
  layout: "fitColumns", //fit columns to width of table (optional)
  pagination: "local", //enable local pagination.
  paginationSize: 10, // this option can take any positive integer value
  ajaxURL: "/api/user-beer-preferences",
  columns: [ //Define Table Columns
    {title: "ID", field: "id"},
    {title: "Name", field: "name"},
    {
      title: "Drunk Beer Before",
      field: "drunk_before",
      formatter: "tickCross",
      editor: true
    },
    {
      title: "Got Drunk",
      field: "got_drunk",
      formatter: "tickCross",
      editor: true
    },
    {
      title: "Last Time",
      field: "last_time",
      editor: dateEditor,
      formatter: "datetime",
      formatterParams: {
        outputFormat: "YYYY-MM-DD",
        invalidPlaceholder: "N/A"
      }
    },
    {title: "Rating", field: "rating", formatter: "star", editor: true},
    {title: "Comment", field: "comment", editor: "input"}
  ],
  dataChanged: uploadData
});