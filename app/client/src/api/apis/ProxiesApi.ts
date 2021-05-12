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


import * as runtime from '../runtime';
import {
    ErrorMessage,
    ErrorMessageFromJSON,
    ErrorMessageToJSON,
    ProxyId,
    ProxyIdFromJSON,
    ProxyIdToJSON,
    ProxySpec,
    ProxySpecFromJSON,
    ProxySpecToJSON,
} from '../models';

export interface ProxyCreateRequest {
    proxySpec: ProxySpec;
}

export interface ProxyRemoveRequest {
    proxyId: string;
}

/**
 * 
 */
export class ProxiesApi extends runtime.BaseAPI {

    /**
     * List of all proxies.
     */
    async proxiesListRaw(): Promise<runtime.ApiResponse<Array<ProxySpec>>> {
        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/proxies`,
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        });

        return new runtime.JSONApiResponse(response, (jsonValue) => jsonValue.map(ProxySpecFromJSON));
    }

    /**
     * List of all proxies.
     */
    async proxiesList(): Promise<Array<ProxySpec>> {
        const response = await this.proxiesListRaw();
        return await response.value();
    }

    /**
     * Create new proxy.
     */
    async proxyCreateRaw(requestParameters: ProxyCreateRequest): Promise<runtime.ApiResponse<ProxyId>> {
        if (requestParameters.proxySpec === null || requestParameters.proxySpec === undefined) {
            throw new runtime.RequiredError('proxySpec','Required parameter requestParameters.proxySpec was null or undefined when calling proxyCreate.');
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        const response = await this.request({
            path: `/proxies`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: ProxySpecToJSON(requestParameters.proxySpec),
        });

        return new runtime.JSONApiResponse(response, (jsonValue) => ProxyIdFromJSON(jsonValue));
    }

    /**
     * Create new proxy.
     */
    async proxyCreate(proxySpec: ProxySpec): Promise<ProxyId> {
        const response = await this.proxyCreateRaw({ proxySpec: proxySpec });
        return await response.value();
    }

    /**
     * Remove proxy.
     */
    async proxyRemoveRaw(requestParameters: ProxyRemoveRequest): Promise<runtime.ApiResponse<void>> {
        if (requestParameters.proxyId === null || requestParameters.proxyId === undefined) {
            throw new runtime.RequiredError('proxyId','Required parameter requestParameters.proxyId was null or undefined when calling proxyRemove.');
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/proxies/{proxyId}`.replace(`{${"proxyId"}}`, encodeURIComponent(String(requestParameters.proxyId))),
            method: 'DELETE',
            headers: headerParameters,
            query: queryParameters,
        });

        return new runtime.VoidApiResponse(response);
    }

    /**
     * Remove proxy.
     */
    async proxyRemove(proxyId: string): Promise<void> {
        await this.proxyRemoveRaw({ proxyId: proxyId });
    }

}
