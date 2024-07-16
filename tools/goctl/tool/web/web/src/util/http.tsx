import axios from "axios";

const Paths = {
  ParseBodyPath: "/api/request/body/parse",
  APIBuildPath: "/api/generate",
};

export type ParseBodyResult = {
  form: ParseBodyForm[];
};

export type APIBuildResult = {
  api: string;
};

export type ParseBodyForm = {
  name: string;
  type: string;
  optional?: boolean;
  defaultValue?: string;
  checkEnum?: boolean;
  enumValue?: string;
  lowerBound?: number;
  upperBound?: number;
};

function postJSON<T>(
  path: string,
  data: {},
  callback: (data: T) => void,
  catchError: (err: string) => void,
): void {
  axios
    .post(path, data)
    .then(function (response) {
      if (response.status === 200) {
        let data = response.data;
        if (data.code === 0) {
          callback(data.data);
        } else {
          catchError(data.msg);
        }
      }
    })
    .catch((err) => {
      console.log(err);
      catchError(err.toString());
    });
}

export const Http = {
  ParseBodyFromJson: (
    json: string,
    callback: (data: ParseBodyForm[]) => void,
    catchError: (err: string) => void,
  ) => {
    postJSON<ParseBodyResult>(
      Paths.ParseBodyPath,
      {
        json: json,
      },
      (data) => {
        callback(data.form);
      },
      catchError,
    );
  },
  Build: (
    param: any,
    callback: (data: string) => void,
    catchError: (err: string) => void,
  ) => {
    postJSON<APIBuildResult>(
      Paths.APIBuildPath,
      param,
      (data) => {
        callback(data.api);
      },
      catchError,
    );
  },
};
