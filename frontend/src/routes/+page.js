/** @type {import('./$types').PageLoad} */
import {getRandomSentence} from "$lib/utils";
import { version } from '../../package.json';


export async function load() {
    const randomSentence = getRandomSentence();

    return {
        randomSentence: randomSentence,
        version: version
    };
};