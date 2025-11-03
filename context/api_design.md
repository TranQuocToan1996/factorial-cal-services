# Describe the idea of the project

# Request and response rules
- Use RESTful API

API 1: Request calculate number:
POST /factorial

body request {
    "number": "10"
}

body response on success with calulated res {
    "code": 200,
    "status": "ok",
    "message": "done",
    "data": {
        "number": "10",
        "factorial_result": "3628800"
    }
}

body response on success with wait res {
    "code": 200,
    "status": "ok",
    "message": "calculating",
    "data": {}
}

body response on fail {
    "code": 500,
    "status": "fail",
    "message": "some error",
    "data": {}
}

---------------

API 2: get calculate number:
GET /factorial?number=10

body response on success with calulated res {
    "code": 200,
    "status": "ok",
    "message": "done",
    "data": {
        "number": "10",
        "factorial_result": "3628800"
    }
}

body response on success with wait res {
    "code": 200,
    "status": "ok",
    "message": "calculating",
    "data": {}
}

body response on exceed threshold number {
    "code": 200,
    "status": "exceed_threshold",
    "message": "number exceed threshold",
    "data": {}
}

body response on fail {
    "code": 500,
    "status": "fail",
    "message": "some error",
    "data": {}
}


API 3: Get metadata
GET /factorial/metadata?number=10000000000

body response on success with calulated res {
    "code": 200,
    "status": "ok",
    "message": "done",
    "data": {
        "id": "uuid-v4",
        "number": "10",
        "factorial_result": "3628800",
        "s3_key": "s3_key",
        "checksum": "checksum",
        "status": "done",
        "created_at": "created_at",
        "updated_at": "updated_at"
    }
}

body response on success with wait res {
    "code": 200,
    "status": "ok",
    "message": "calculating",
    "data": {}
}

body response on exceed threshold number {
    "code": 200,
    "status": "exceed_threshold",
    "message": "number exceed threshold",
    "data": {}
}

body response on fail {
    "code": 500,
    "status": "fail",
    "message": "some error",
    "data": {}
}
