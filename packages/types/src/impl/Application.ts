import type { DiscoverSectionItem } from "../DiscoverSectionItem"
import type { DiscoverSection } from "../HomeSection"
import type { PagedResults } from "../PagedResults"
import type { SearchFilter } from "../SearchFilter"
import { Request } from "../Request"
import { Response } from "../Response"

export {};

type RequestInterceptor = (request: Request) => Promise<Request>;
type ResponseInterceptor = (
  request: Request,
  response: Response,
  data: ArrayBuffer,
) => Promise<ArrayBuffer>;

/**
 * @param proposedRequest The `Request` to the new location specified by the redirect response.
 * @param redirectedResponse The `Response` containing the server's response to the original request.
 * @returns Return the proposed request or a modified request to follow the redirect, or undefined to cancel the redirect
 */
type RedirectHandler = (
  proposedRequest: Request,
  redirectedResponse: Response,
) => Promise<Request | undefined>;

declare global {
  namespace Application {
    // Global
    const isResourceLimited: boolean;
    function decodeHTMLEntities(str: string): string;
    function sleep(seconds: number): Promise<void>;

    // Discover Section
    /**
     * @description If sections are registered using this method the app will not call {@link DiscoverSectionProviding.getDiscoverSections} unless {@link Application.invalidateDiscoverSections} is called.
     * @param selector if provided, the app will call the method, otherwise defaults to {@link DiscoverSectionProviding.getDiscoverSectionItems}
     * @deprecated register sections in {@link DiscoverSectionProviding.getDiscoverSections} by implementing {@link DiscoverSectionProviding}
     */
    function registerDiscoverSection(
      section: DiscoverSection,
      selector?: SelectorID<
        (
          section: DiscoverSection,
          metadata: unknown | undefined,
        ) => Promise<PagedResults<DiscoverSectionItem>>
      >,
    ): void;
    function unregisterDiscoverSection(sectionId: string): void;
    function registeredDiscoverSections(): DiscoverSection[];
    /** Invalidate discover section cache (removes all discover sections) */
    function invalidateDiscoverSections(): void;

    // Request Manager
    function registerInterceptor(
      interceptorId: string,
      interceptRequestSelectorId: SelectorID<RequestInterceptor>,
      interceptResponseSelectorId: SelectorID<ResponseInterceptor>,
    ): void;
    function unregisterInterceptor(interceptorId: string): void;
    function setRedirectHandler(
      redirectHandlerSelectorId: SelectorID<RedirectHandler>,
    ): void;
    function getDefaultUserAgent(): Promise<string>;
    function scheduleRequest(
      request: Request,
    ): Promise<[Response, ArrayBuffer]>;

    // Raw Data
    function arrayBufferToUTF8String(arrayBuffer: ArrayBuffer): string;
    function arrayBufferToASCIIString(arrayBuffer: ArrayBuffer): string;
    function arrayBufferToUTF16String(arrayBuffer: ArrayBuffer): string;

    // Search Filters
    function registerSearchFilter(searchFilter: SearchFilter): void;
    function unregisterSearchFilter(id: string): void;
    function registeredSearchFilters(): SearchFilter[];

    // State Manager
    function getSecureState(key: string): unknown | undefined;
    function setSecureState(value: unknown, key: string): void;
    function getState(key: string): unknown | undefined;
    function setState(value: unknown, key: string): void;
    /**
     * Clears all saved state.
     *
     * *Note: Does not clear secure state.*
     */
    function resetAllState(): void;
  }
}