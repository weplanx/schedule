<?php

namespace lumen\extra\middleware;

use Closure;
use Illuminate\Http\Request;
use lumen\extra\facade\Auth;

class AuthVerify
{
    protected $scene;

    /**
     * Handle an incoming request.
     *
     * @param \Illuminate\Http\Request $request
     * @param \Closure $next
     * @return mixed
     */
    public function handle(Request $request, Closure $next)
    {
        if (empty($this->scene)) {
            return $next($request);
        }

        $result = Auth::verify($this->scene);
        if (!$result) {
            return response()->json([
                'error' => 1,
                'msg' => 'token invalid'
            ]);
        } else {
            return $next($request);
        }
    }
}
