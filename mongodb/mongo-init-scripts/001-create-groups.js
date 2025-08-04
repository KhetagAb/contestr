db.createCollection("groups", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            required: ["school", "year", "parallel", "students", "judgeSystem", "contests"],
            additionalProperties: false,
            properties: {
                school: {
                    enum: ["ITMO", "LKSH", "Sirius"],
                },
                year: {
                    bsonType: "int",
                    minimum: 2000,
                    maximum: 2030,
                },
                shift: {
                    bsonType: "string",
                    description: "suffix for year, for example 'august'",
                },
                parallel: {
                    bsonType: "",
                },
                students: {
                    bsonType: "array",
                    description: "student's objectId documents",
                    items: {
                        bsonType: "object",
                        required: ["handle", "studentId"],
                        properties: {
                            handle: {
                                bsonType: "string"
                            },
                            studentId: {
                                bsonType: "objectId"
                            }
                        }
                    }
                },
                judgeSystem: {
                    enum: ["codeforces", "ejudge"]
                },
                contests: {
                    bsonType: "array",
                    description: "contests identifiers",
                    items: {
                        bsonType: "int",
                    }
                }
            }
        }
    }
});
