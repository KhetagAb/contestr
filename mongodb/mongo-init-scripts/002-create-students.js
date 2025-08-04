db.createCollection("students", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            required: ["name", "judgeHandle"],
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
                        bsonType: "objectId"
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
