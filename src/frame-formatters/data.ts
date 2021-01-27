import { Observable } from 'rxjs';
import { switchMap as switchMap$ } from 'rxjs/operators';
import { DataFrame, DataQueryResponse } from '@grafana/data';

/**
 * DataFrameFormatter
 */
export class DataFrameFormatter {
  /**
   * update
   * @param request
   */
  async update(request: Observable<DataQueryResponse>): Promise<DataFrame> {
    return request.pipe(switchMap$((response) => response.data)).toPromise();
  }
}
