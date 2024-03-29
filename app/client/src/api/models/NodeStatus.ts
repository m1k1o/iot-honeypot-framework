/* tslint:disable */
/* eslint-disable */
/**
 * IoT Honeypot API
 * IoT Honeypot API 
 *
 * The version of the OpenAPI document: 1.0.0
 * Contact: xsedivy@stuba.sk
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import { exists, mapValues } from '../runtime';
/**
 * 
 * @export
 * @interface NodeStatus
 */
export interface NodeStatus {
    /**
     * 
     * @type {string}
     * @memberof NodeStatus
     */
    addr?: string;
    /**
     * 
     * @type {string}
     * @memberof NodeStatus
     */
    state?: NodeStatusStateEnum;
}

/**
* @export
* @enum {string}
*/
export enum NodeStatusStateEnum {
    Unknown = 'unknown',
    Down = 'down',
    Ready = 'ready',
    Disconnected = 'disconnected'
}

export function NodeStatusFromJSON(json: any): NodeStatus {
    return NodeStatusFromJSONTyped(json, false);
}

export function NodeStatusFromJSONTyped(json: any, ignoreDiscriminator: boolean): NodeStatus {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'addr': !exists(json, 'addr') ? undefined : json['addr'],
        'state': !exists(json, 'state') ? undefined : json['state'],
    };
}

export function NodeStatusToJSON(value?: NodeStatus | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'addr': value.addr,
        'state': value.state,
    };
}


