let data = require('./data.json')
let filtered = require("./filtered.json")

// let doubled = filtered.map(e => {
//     return e.exports.map(el => {
//         return el.properties.filter(p => {
//             return p.type === "Array of Structs"
//         })
//     })
// }).flat(9999)
let doubled = filtered.filter(e => {
    return e.exports.filter(exp => {
        return exp.name.startsWith("Default__")
    }).length < 1
})
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