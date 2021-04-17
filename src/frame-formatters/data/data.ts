import { Observable } from 'rxjs';
import { switchMap as switchMap$ } from 'rxjs/operators';
import { DataFrame, DataQueryResponse } from '@grafana/data';

/**
 * Data Frame Formatter
 */
export class DataFrameFormatter {
  /**
   * Update
   *
   * @param request
   */
  async update(request: Observable<DataQueryResponse>): Promise<DataFrame> {
    return request.pipe(switchMap$((response) => response.data)).toPromise();
  }
}
