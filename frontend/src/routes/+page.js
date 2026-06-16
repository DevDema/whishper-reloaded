/** @type {import('./$types').PageLoad} */
import {getRandomTaglineIndex} from "$lib/utils";
import { version } from '../../package.json';


export async function load() {
    const taglineIndex = getRandomTaglineIndex();

    return {
        taglineIndex: taglineIndex,
        version: version
    };
};