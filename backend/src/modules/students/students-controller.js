const asyncHandler = require("express-async-handler");
const { getAllStudents, addNewStudent, getStudentDetail, setStudentStatus, updateStudent } = require("./students-service");

const handleGetAllStudents = asyncHandler(async (req, res) => {
    //write your code
    const students = await getAllStudents(req.params);
    res.json({ students });
});

const handleAddStudent = asyncHandler(async (req, res) => {
    //write your code
    const message =  await addNewStudent(req.body);
    res.json(message);
});

const handleUpdateStudent = asyncHandler(async (req, res) => {
    //write your code

});

const handleGetStudentDetail = asyncHandler(async (req, res) => {
    //write your code
    const { id } = req.params;
    const student = await getStudentDetail(id);
    res.json(student);
});

const handleStudentStatus = asyncHandler(async (req, res) => {
    //write your code

});

module.exports = {
    handleGetAllStudents,
    handleGetStudentDetail,
    handleAddStudent,
    handleStudentStatus,
    handleUpdateStudent,
};
