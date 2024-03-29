import * as fs from "fs"

import chain from 'stream-chain';
import parser from 'stream-json';
import streamValues from 'stream-json/streamers/StreamValues.js';

const pipeline = chain.chain([
    fs.createReadStream('data-build.json'),
    parser.parser(),
    streamValues.streamValues(),
]);

const promise = new Promise((resolve, reject) => {
    pipeline.on("end", resolve)
})

let data_build = {}

pipeline.on("data", newData => {
    data_build = newData.value
})

await promise

import data from "./data-build.json" assert { type: "json" };
import one from "./one-pretty.json" assert { type: "json" };
import filtered from "./filtered.json" assert { type: "json" };

// let build = data.filter(e => {
//     return e.export_record.file_name.includes("Build_")
// })
// fs.writeFileSync("data-build.json", JSON.stringify(build))

// let conveyors = data.filter(e => {
//     return e.export_record.file_name.toLowerCase().includes("conveyor")
// })
//
// let classProps = conveyors.filter(e => {
//     return e.exports[0].data.properties.length > 0
// })

debugger
// let classExports = filtered.map(e => {
//     return e.exports.filter(el => {
//         return !el.object_name.startsWith("Default__") && el.object_name.endsWith("_C")
//     })
// }).flat(2)
// let classesWithNewShit = classExports.filter(cla => {
//     return cla.properties.filter(p => {
//                 let valid = true
//                     const names = ["Timelines","Animations", "SimpleConstructionScript", "WidgetTree", "bClassRequiresNativeTick", "Bindings", "DynamicBindingObjects", "UberGraphFunction", "InheritableComponentHandler"]
//                         names.forEach(s => {
//                         if (p.name === s) {
//                             valid = false
//                         }
//                 })
//                 return valid
//             }).length > 0
// })
// let builds = filtered.filter(e => {
//     return e.filename.includes("Build_") && e.exports.length !== 7
// })
// let doubled = filtered.map(e => {
//     return e.exports.map(el => {
//         if (!el.object_name.startsWith("Default__") && el.object_name.endsWith("_C")) {
//             return el.properties.filter(p => {
//                 let valid = true
//                     const names = ["SimpleConstructionScript", "WidgetTree", "bClassRequiresNativeTick", "Bindings", "DynamicBindingObjects", "UberGraphFunction", "InheritableComponentHandler"]
//                         names.forEach(s => {
//                         if (p.name === s) {
//                             valid = false
//                         }
//                 })
//                 return valid
//             })
//         }
//     })
// })
// let doubled = filtered.filter(e => {
//     return e.exports.filter(exp => {
//         return exp.name.startsWith("Default__")
//     }).length < 1
// })
// let doubled = data.map(e => {
//     return e.summary.imports.map(el => {
//         r = []
//         i = (el.outer_index+1) * -1
//         outer = e.summary.imports[i]
//         if (el.outer_index !== 0 && outer.class_name !== "Package") {
//             r.push(el)
//             r.push(outer)
//             outer2 = e.summary.imports[(outer.outer_index+1) * -1]
//             r.push(outer2)
//             if (outer2.outer_index !== 0) {
//                 r.push(e.summary.imports[(outer2.outer_index+1) * -1])
//             }
//         }
//         return r
//     }).filter(e => {
//         return e.length > 3
//     })
// }).flat(1)
debugger

function isFilteredStruct(e) {
    e.type === "Struct" || e.type === "Array of Structs" || e.type === "Map of Objects to Structs"
}

const plist = [
    "ObjectProperty",
    "TextProperty",
    "EnumProperty",
    "ByteProperty",
    "NameProperty",
    "FloatProperty",
    "ArrayProperty",
    "IntProperty"
]

function isObject(val) {
    return val && typeof val === 'object'
}

function deepForEach(obj, func) {
    Object.keys(obj).forEach(key => {
        func(obj, key)
        if (isObject(obj[key])) {
            deepForEach(obj[key], func)
        }
    })
}

function isArrayOfStructs(e) {
    return isArrayOfType(e, "Struct")
}

function isArrayOfType(e, t) {
    return isArray(e)
        && isOfTypeStr(e.tag_data, t)
}

function isMapWithOneOfType(e, t) {
    return isMap(e)
        && (isOfTypeStr(e.tag_data.key_type, t) || isOfTypeStr(e.tag_data.value_type, t))
}

function isMap(e) {
    return isOfType(e, "Map")
}

function isArray(e) {
    return isOfType(e, "Array")
}

function isStruct(e) {
    return isOfType(e, "Struct")
}

function isOfType(e, t) {
    return isOfTypeStr(e.property_type, t)
}

function isOfTypeStr(s, t) {
    return clean(s) === t + "Property"
}

function clean(s) {
    return s.replace("\u0000", "")
}